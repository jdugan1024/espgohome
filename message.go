package espgohome

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

// MessageID is an enum of all the known messages types
type MessageID uint64

//go:generate stringer -type=MessageID

const (
	HelloRequestID                          MessageID = 1
	HelloResponseID                         MessageID = 2
	ConnectRequestID                        MessageID = 3
	ConnectResponseID                       MessageID = 4
	DisconnectRequestID                     MessageID = 5
	DisconnectResponseID                    MessageID = 6
	PingRequestID                           MessageID = 7
	PingResponseID                          MessageID = 8
	DeviceInfoRequestID                     MessageID = 9
	DeviceInfoResponseID                    MessageID = 10
	ListEntitiesRequestID                   MessageID = 11
	ListEntitiesDoneResponseID              MessageID = 19
	SubscribeStatesRequestID                MessageID = 20
	ListEntitiesBinarySensorResponseID      MessageID = 12
	BinarySensorStateResponseID             MessageID = 21
	ListEntitiesCoverResponseID             MessageID = 13
	CoverStateResponseID                    MessageID = 22
	CoverCommandRequestID                   MessageID = 30
	ListEntitiesFanResponseID               MessageID = 14
	FanStateResponseID                      MessageID = 23
	FanCommandRequestID                     MessageID = 31
	ListEntitiesLightResponseID             MessageID = 15
	LightStateResponseID                    MessageID = 24
	LightCommandRequestID                   MessageID = 32
	ListEntitiesSensorResponseID            MessageID = 16
	SensorStateResponseID                   MessageID = 25
	ListEntitiesSwitchResponseID            MessageID = 17
	SwitchStateResponseID                   MessageID = 26
	SwitchCommandRequestID                  MessageID = 33
	ListEntitiesTextSensorResponseID        MessageID = 18
	TextSensorStateResponseID               MessageID = 27
	SubscribeLogsRequestID                  MessageID = 28
	SubscribeLogsResponseID                 MessageID = 29
	SubscribeHomeassistantServicesRequestID MessageID = 34
	HomeassistantServiceResponseID          MessageID = 35
	SubscribeHomeAssistantStatesRequestID   MessageID = 38
	SubscribeHomeAssistantStateResponseID   MessageID = 39
	HomeAssistantStateResponseID            MessageID = 40
	GetTimeRequestID                        MessageID = 36
	GetTimeResponseID                       MessageID = 37
	ListEntitiesServicesResponseID          MessageID = 41
	ExecuteServiceRequestID                 MessageID = 42
	ListEntitiesCameraResponseID            MessageID = 43
	CameraImageResponseID                   MessageID = 44
	CameraImageRequestID                    MessageID = 45
	ListEntitiesClimateResponseID           MessageID = 46
	ClimateStateResponseID                  MessageID = 47
	ClimateCommandRequestID                 MessageID = 48
)

func decodeMessage(raw []byte, msgType MessageID) (proto.Message, error) {
	switch msgType {
	case ListEntitiesBinarySensorResponseID:
		resp := &ListEntitiesBinarySensorResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case CameraImageRequestID:
		resp := &CameraImageRequest{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case SubscribeLogsResponseID:
		resp := &SubscribeLogsResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case SubscribeLogsRequestID:
		resp := &SubscribeLogsRequest{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case SensorStateResponseID:
		resp := &SensorStateResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case CoverStateResponseID:
		resp := &CoverStateResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case ListEntitiesLightResponseID:
		resp := &ListEntitiesLightResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case ListEntitiesDoneResponseID:
		resp := &ListEntitiesDoneResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case GetTimeRequestID:
		resp := &GetTimeRequest{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case ListEntitiesTextSensorResponseID:
		resp := &ListEntitiesTextSensorResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case DeviceInfoRequestID:
		resp := &DeviceInfoRequest{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case ConnectRequestID:
		resp := &ConnectRequest{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case TextSensorStateResponseID:
		resp := &TextSensorStateResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case SubscribeStatesRequestID:
		resp := &SubscribeStatesRequest{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case HomeAssistantStateResponseID:
		resp := &HomeAssistantStateResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case DisconnectResponseID:
		resp := &DisconnectResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case ListEntitiesSensorResponseID:
		resp := &ListEntitiesSensorResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case ListEntitiesCameraResponseID:
		resp := &ListEntitiesCameraResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case BinarySensorStateResponseID:
		resp := &BinarySensorStateResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case SubscribeHomeAssistantStateResponseID:
		resp := &SubscribeHomeAssistantStateResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case ClimateStateResponseID:
		resp := &ClimateStateResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case CameraImageResponseID:
		resp := &CameraImageResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case ListEntitiesCoverResponseID:
		resp := &ListEntitiesCoverResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case GetTimeResponseID:
		resp := &GetTimeResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case SwitchCommandRequestID:
		resp := &SwitchCommandRequest{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case LightStateResponseID:
		resp := &LightStateResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case LightCommandRequestID:
		resp := &LightCommandRequest{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case DeviceInfoResponseID:
		resp := &DeviceInfoResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case HelloResponseID:
		resp := &HelloResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case FanCommandRequestID:
		resp := &FanCommandRequest{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case PingRequestID:
		resp := &PingRequest{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case SubscribeHomeassistantServicesRequestID:
		resp := &SubscribeHomeassistantServicesRequest{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case ExecuteServiceRequestID:
		resp := &ExecuteServiceRequest{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case CoverCommandRequestID:
		resp := &CoverCommandRequest{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case ListEntitiesServicesResponseID:
		resp := &ListEntitiesServicesResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case DisconnectRequestID:
		resp := &DisconnectRequest{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case ListEntitiesSwitchResponseID:
		resp := &ListEntitiesSwitchResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case PingResponseID:
		resp := &PingResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case ConnectResponseID:
		resp := &ConnectResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case SubscribeHomeAssistantStatesRequestID:
		resp := &SubscribeHomeAssistantStatesRequest{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case FanStateResponseID:
		resp := &FanStateResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case HelloRequestID:
		resp := &HelloRequest{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case SwitchStateResponseID:
		resp := &SwitchStateResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case ClimateCommandRequestID:
		resp := &ClimateCommandRequest{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case ListEntitiesFanResponseID:
		resp := &ListEntitiesFanResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case ListEntitiesRequestID:
		resp := &ListEntitiesRequest{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case ListEntitiesClimateResponseID:
		resp := &ListEntitiesClimateResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case HomeassistantServiceResponseID:
		resp := &HomeassistantServiceResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	default:
		err := fmt.Errorf("unsupported message: %d", msgType)
		return nil, err
	}
}
