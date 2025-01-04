package radio

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/lmatte7/gomesh"
	"github.com/lmatte7/gomesh/github.com/meshtastic/gomeshproto"
)

type userPos struct {
	User *gomeshproto.User
	Pos  *gomeshproto.Position
}

type RadioClient struct {
	radio   *gomesh.Radio
	userdir map[string]userPos
}

type BroadcastOpt struct {
	To      int64
	Channel int64
}

func NewRadioClient(opts ...func(*RadioClient)) *RadioClient {
	c := &RadioClient{
		userdir: map[string]userPos{},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

func WithPort(port string) func(*RadioClient) {
	return func(r *RadioClient) {
		//not returning an errors since it'll panic here anyway
		r.radio = &gomesh.Radio{}
		if err := r.radio.Init(port); err != nil {
			panic(err)
		}
	}
}

func (r *RadioClient) Close() {
	r.radio.Close()
}

func (r *RadioClient) UserDir() map[string]userPos {
	return r.userdir
}

func (r *RadioClient) Info() error {
	nodes, err := r.radio.GetRadioInfo()
	if err != nil {
		return err
	}
	for _, node := range nodes {
		if ni := node.GetNodeInfo(); ni == nil {
			continue
		}
		if user := node.GetNodeInfo().GetUser(); user != nil {
			up := userPos{
				User: user,
				Pos:  node.GetNodeInfo().GetPosition(),
			}
			r.userdir[up.User.Id] = up
		}
	}
	return nil
}

func (r *RadioClient) Send(msg string, to, channel int64) {
	r.radio.SendTextMessage(msg, to, channel)
}

func (r *RadioClient) UserPos(nodeID string) (*gomeshproto.User, *gomeshproto.Position, error) {
	if err := r.Info(); err != nil {
		return nil, nil, err
	}
	if up, ok := r.userdir[nodeID]; ok {
		return up.User, up.Pos, nil
	}
	return nil, nil, nil
}

func (r *RadioClient) Receive(c chan<- []byte, done chan<- bool) {
	defer close(c) // Ensure the channel is closed when this function exits
	for {
		packets, err := r.radio.ReadResponse(true)
		if err != nil {
			// Signal done and exit
			log.Panic(err)
			done <- true
			return
		}

		for _, packet := range packets {
			if p := packet.GetPacket(); p != nil {
				if d := p.GetDecoded(); d != nil {
					if slices.Contains([]string{
						"TELEMETRY_APP", "NODEINFO_APP",
						"POSITION_APP", "ROUTING_APP",
					}, string(d.Portnum.String())) {
						continue
					}
					switch d.Portnum.String() {
					case "TEXT_MESSAGE_APP":
						encoded, _ := json.Marshal(map[string]interface{}{
							"from":    p.From,
							"to":      d.Source,
							"data":    packet,
							"portnum": d.Portnum.String(),
							"payload": string(d.Payload),
						})
						fmt.Println(string(encoded))
						c <- encoded
					}

				}
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
}
