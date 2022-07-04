package notificationHandler

import "tsn-service/pkg/structures/configuration"

// TODO: Implement correctly with actual values that should be used
// from the reqData or whatever is required to fill it.

// TODO: Structure doesn't match response.json... FIX IT
func generateBaseResponse(reqData *configuration.Request) (*configuration.ConfigResponse, error) {
	var baseResp = &configuration.ConfigResponse{
		Version: 123, // random value for now
		Responses: []*configuration.Response{
			{
				StatusGroup: &configuration.StatusGroup{
					StrId: &configuration.StreamId{
						MacAddress: "11:22:33:44:55", // random value for now
						UniqueId:   "123",            // random value for now
					},
					StatusInfo: &configuration.StatusInfo{
						TalkerStatus:   1,   // random value for now
						ListenerStatus: 1,   // random value for now
						FailureCode:    123, // random value for now
					},
					FailedInterfaces: []*configuration.InterfaceId{
						{
							MacAddress:    "11:22:33:44:55", // random value for now
							InterfaceName: "I-do-not-know",  // random value for now
						},
					},
					StatusTalkerListener: []*configuration.TalkerListenerStatus{
						{
							AccumulatedLatency: &configuration.AccumulatedLatency{
								AccumulatedLatency: 123, // random value for now
							},
							InterfaceConfiguration: []*configuration.InterfaceConfiguration{
								{
									InterfaceId: &configuration.InterfaceId{
										MacAddress:    "11:22:33:44:55", // random value for now
										InterfaceName: "test",           // random value for now
									},
									Type: 123, // random value for now
									MacAddr: &configuration.IeeeMacAddress{
										DestinationMac: "destMac", // random value for now
										SourceMac:      "srcMac",  // random value for now
									},
									VlanTag: &configuration.IeeeVlanTag{
										PriorityCodePoint: 1, // random value for now
										VlanId:            1, // random value for now
									},
									Ipv4Tup: &configuration.Ipv4Tuple{},
									Ipv6Tup: &configuration.Ipv6Tuple{},
									TimeAwareOffset: &configuration.TimeAwareOffset{
										Offset: 123, // random value for now
									},
								},
							},
						},
					},
					EndStationInterfaces: []*configuration.Interface{
						{
							Index: 0, // random value for now
							InterfaceId: &configuration.InterfaceId{
								MacAddress:    "11:22:33:44:55", // random value for now
								InterfaceName: "testInterface",  // random value for now
							},
						},
					},
				},
			},
		},
	}

	return baseResp, nil
}
