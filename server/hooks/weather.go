package hooks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"

	"mesh-mqtt/clients/noaa"
	"mesh-mqtt/clients/radio"

	"mesh-mqtt/server/models"

	"github.com/cockroachdb/pebble"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/system"
)

type WeatherHookOptions struct {
	Server *mqtt.Server
	Radio  *radio.RadioClient
	NOAA   *noaa.NOAAClient
	DB     *pebble.DB
}

type WeatherHook struct {
	mqtt.HookBase
	config     *WeatherHookOptions
	nextStatus int64
}

func (r *WeatherHook) ID() string {
	return "weather"
}

func (r *WeatherHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnSysInfoTick,
	}, []byte{b})
}

func (r *WeatherHook) Init(config any) error {
	if _, ok := config.(*WeatherHookOptions); !ok && config != nil {
		return mqtt.ErrInvalidConfigType
	}

	r.config = config.(*WeatherHookOptions)
	check := []interface{}{
		r.config.Server,
		r.config.Radio,
		r.config.NOAA,
		r.config.DB,
	}
	if slices.ContainsFunc(check, func(arg any) bool {
		return arg == nil
	}) {
		return mqtt.ErrInvalidConfigType
	}

	// wait for node to come online
	r.nextStatus = time.Now().Add(60 * time.Second).Unix()
	return nil
}

func (r *WeatherHook) OnSysInfoTick(si *system.Info) {
	if time.Now().Unix() <= r.nextStatus {
		return
	}

	// send the status update to the system channel
	r.config.Radio.Send(fmt.Sprintf("t:%v, up: %v(s), clients:%v, memalloc: %v",
		time.Now().Unix(),
		si.Uptime,
		si.ClientsConnected,
		si.MemoryAlloc),
		0, 2)
	// top of the hour, every hour
	r.nextStatus = time.Now().Truncate(time.Hour).Add(time.Hour).Unix()

	// query the radio to determine which nodes expose gps coords to fetch weather forecasts
	r.config.Radio.Info()
	// outer iterator to keep from querying noaa twice for the same cwa/grid
	dedupe := map[string]bool{}
	for key, userPos := range r.config.Radio.UserDir() {
		if userPos.Pos == nil || key == "" {
			continue
		}
		nr := models.NodeRecord{}
		// filter any nodes not connected to the mqtt instance
		nodebytes, closer, err := r.config.DB.Get([]byte(key))
		if err != nil {
			if errors.Is(err, pebble.ErrNotFound) {
				r.config.Server.Log.Info("error fetching client not found: ", key)
			}
			continue
		} else {
			// only close a pebble db query if its successful
			closer.Close()
		}

		if len(nodebytes) == 0 {
			continue
		}

		if err := json.Unmarshal(nodebytes, &nr); err != nil {
			r.config.Server.Log.Error("error during unmarshal client ", key, err)
			continue
		}

		// fetch noaa cwa and grid data for forecasts
		if nr.GridID == "" {
			lat, long := float64(userPos.Pos.LatitudeI)/1e7, float64(userPos.Pos.LongitudeI)/1e7
			points, err := r.config.NOAA.Points(lat, long)
			if err != nil {
				r.config.Server.Log.Error("error fetching points from noaa @", lat, ",", long)
				continue
			}
			nr.GridID = points.Properties.GridID
			nr.GridX = points.Properties.GridX
			nr.GridY = points.Properties.GridY

			// cache it so we don't have to query noaa twice for the same node
			if nodebytes, err = json.Marshal(nr); err != nil {
				r.config.Server.Log.Error("error update user ", key)
				continue
			}
			r.config.DB.Set([]byte(key), nodebytes, nil)

		}
		// if we haven't come across this grid before fetch the forecast and send the forecast for this hour
		if _, ok := dedupe[fmt.Sprintf("%s%v%v", nr.GridID, nr.GridX, nr.GridY)]; !ok {
			forecast, err := r.config.NOAA.GridPointsForecast(nr.GridID, int64(nr.GridX), int64(nr.GridY))
			if err != nil {
				r.config.Server.Log.Error("unable to get gridpoints ", key, " ", err)
				return
			}

			if len(forecast.Properties.Periods) == 0 {
				r.config.Server.Log.Error("invalid forecast length ", key)
				continue
			}
			msg := fmt.Sprintf("mesh hourly forecast: %s precip: %v %% chance.",
				forecast.Properties.Periods[0].DetailedForecast,
				forecast.Properties.Periods[0].ProbabilityOfPrecipitation.Value,
			)
			// todo: configure channels for weather updates
			// for now just send them to the private system channel as to not be annoying
			r.config.Radio.Send(msg, 0, 2)
			dedupe[fmt.Sprintf("%s%v%v", nr.GridID, nr.GridX, nr.GridY)] = true
		}
	}
}
