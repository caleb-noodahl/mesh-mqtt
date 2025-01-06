package hooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mesh-mqtt/clients/radio"
	"mesh-mqtt/server/models"
	"slices"
	"time"

	"github.com/cockroachdb/pebble"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
)

type UserHookOptions struct {
	Server *mqtt.Server
	Radio  *radio.RadioClient
	DB     *pebble.DB
}

type UserHook struct {
	mqtt.HookBase
	config *UserHookOptions
}

func (h *UserHook) Init(config any) error {
	if _, ok := config.(*UserHookOptions); !ok && config != nil {
		return mqtt.ErrInvalidConfigType
	}
	h.config = config.(*UserHookOptions)
	check := []interface{}{
		h.config.Server,
		h.config.Radio,
		h.config.DB,
	}
	if slices.ContainsFunc(check, func(arg any) bool {
		return arg == nil
	}) {
		return mqtt.ErrInvalidConfigType
	}
	return nil
}

func (h *UserHook) ID() string {
	return "user"
}

func (h *UserHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnect,
		//mqtt.OnDisconnect,
	}, []byte{b})
}

func (h *UserHook) OnConnect(cl *mqtt.Client, pk packets.Packet) error {
	h.config.Server.Log.Debug(fmt.Sprintf("client %s connected", cl.ID))
	data, err := json.Marshal(models.NodeRecord{
		ID:      cl.ID,
		Connect: time.Now().Unix(),
		Last:    time.Now().Unix(),
	})

	go func() {
		time.Sleep(15 * time.Second)
		user, pos, err := h.config.Radio.UserPos(cl.ID)
		if err != nil || user == nil || pos == nil {
			h.config.Server.Log.Warn(fmt.Sprintf("unable to fetch user or pos for client:%s", cl.ID))
			return
		}
		h.config.Server.Log.Info(fmt.Sprintf("client: %s, name: %s, position: %v:%v", cl.ID, user.LongName, pos.LatitudeI, pos.LongitudeI))
	}()

	if err != nil {
		return err
	}
	return h.config.DB.Set([]byte(cl.ID), data, nil)
}
