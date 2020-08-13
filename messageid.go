package espgohome

//go:generate stringer -type=MessageID
type MessageID uint64
const (
	HelloRequestID MessageID = 1
	HelloResponseID MessageID = 2
	ConnectRequestID MessageID = 3
	ConnectResponseID MessageID = 4
	DisconnectRequestID MessageID = 5
	DisconnectResponseID MessageID = 6
	PingRequestID MessageID = 7
	PingResponseID MessageID = 8
	DeviceInfoRequestID MessageID = 9
	DeviceInfoResponseID MessageID = 10
	ListEntitiesRequestID MessageID = 11
	ListEntitiesDoneResponseID MessageID = 19
	SubscribeStatesRequestID MessageID = 20
	ListEntitiesBinarySensorResponseID MessageID = 12
	BinarySensorStateResponseID MessageID = 21
	ListEntitiesCoverResponseID MessageID = 13
	CoverStateResponseID MessageID = 22
	CoverCommandRequestID MessageID = 30
	ListEntitiesFanResponseID MessageID = 14
	FanStateResponseID MessageID = 23
	FanCommandRequestID MessageID = 31
	ListEntitiesLightResponseID MessageID = 15
	LightStateResponseID MessageID = 24
	LightCommandRequestID MessageID = 32
	ListEntitiesSensorResponseID MessageID = 16
	SensorStateResponseID MessageID = 25
	ListEntitiesSwitchResponseID MessageID = 17
	SwitchStateResponseID MessageID = 26
	SwitchCommandRequestID MessageID = 33
	ListEntitiesTextSensorResponseID MessageID = 18
	TextSensorStateResponseID MessageID = 27
	SubscribeLogsRequestID MessageID = 28
	SubscribeLogsResponseID MessageID = 29
	SubscribeHomeassistantServicesRequestID MessageID = 34
	HomeassistantServiceResponseID MessageID = 35
	SubscribeHomeAssistantStatesRequestID MessageID = 38
	SubscribeHomeAssistantStateResponseID MessageID = 39
	HomeAssistantStateResponseID MessageID = 40
	GetTimeRequestID MessageID = 36
	GetTimeResponseID MessageID = 37
	ListEntitiesServicesResponseID MessageID = 41
	ExecuteServiceRequestID MessageID = 42
	ListEntitiesCameraResponseID MessageID = 43
	CameraImageResponseID MessageID = 44
	CameraImageRequestID MessageID = 45
	ListEntitiesClimateResponseID MessageID = 46
	ClimateStateResponseID MessageID = 47
	ClimateCommandRequestID MessageID = 48
)
