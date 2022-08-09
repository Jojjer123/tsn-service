package optimizer

import (
	"errors"
	"strings"
	"tsn-service/pkg/structures/schedule"

	"github.com/onosproject/onos-api/go/onos/topo"
	pb "github.com/openconfig/gnmi/proto/gnmi"
)

// Creates configuration set request from given schedule and network topology
func createConfigurationFromSchedule(sched *schedule.Schedule, topology []topo.Object) (*pb.SetRequest, error) {
	var devicePortMap map[string][]string

	devicePortMap, err := findAllPortsOnDevices(topology)
	if err != nil {
		log.Errorf("Failed finding all ports on devices: %v", err)
		return nil, err
	}

	// log.Infof("Map of devices with their ports looks like: %v", devicePortMap)

	req := &pb.SetRequest{}

	// For every object in the topology
	for _, topoObj := range topology {
		// If device is an entity (switch or other network device)
		if topoObj.Type == topo.Object_ENTITY {
			// Type assert object to entity
			entity := topoObj.Obj.(*topo.Object_Entity)

			// If entity is of kind "netconf-device", hard-coded for now, would be better to check against a list
			// derived from a config file or something listing all kinds of devices in the network.
			if entity.Entity.KindID == "netconf-device" {
				log.Info("Topo device of kind \"netconf-device\" found")

				// Add GCL and gating cycle for all port interfaces on device
				for _, port := range devicePortMap[string(topoObj.ID)] {
					// Create initial gate status updates
					statusChange := getStatusChangeElems(port, string(topoObj.ID), len(sched.TrafficClasses))

					// Create GCL entries for a specific port
					gcl := getGclElems(sched, port, string(topoObj.ID))

					// Create gating cycle entry for a specific port
					adminCycle := getAdminCycleTimeElems(sched.GatingCycle, port, string(topoObj.ID))

					// Create extra time information and set config change to true
					timeInfoAndConfigChange := getFinalElems(port, string(topoObj.ID))

					// Add all updates to set request
					req.Update = append(req.Update, statusChange...)
					req.Update = append(req.Update, gcl...)
					req.Update = append(req.Update, adminCycle...)
					req.Update = append(req.Update, timeInfoAndConfigChange...)
				}
			}
		}
	}

	// log.Infof("Request looks like: %v", *req)

	return req, nil
}

// Finds every port on the devices (looks at all relations and finds the "links")
func findAllPortsOnDevices(topology []topo.Object) (map[string][]string, error) {
	var devicePortMap = map[string][]string{}

	// Find all ports on all devices
	for _, topoObj := range topology {
		if topoObj.Type == topo.Object_RELATION {
			// log.Infof("Checking if \"%s\" already exists in map", string(topoObj.Obj.(*topo.Object_Relation).Relation.SrcEntityID))

			var srcDevice = string(topoObj.Obj.(*topo.Object_Relation).Relation.SrcEntityID)
			var dstDevice = string(topoObj.Obj.(*topo.Object_Relation).Relation.TgtEntityID)
			var srcPort string
			var dstPort string

			var ok bool

			// If relation is created by onos-config (not link, not sure what it is) skip that relation
			if strings.Contains(srcDevice, "onos-config") {
				continue
			}

			var adHocAspect = &topo.AdHoc{}
			if err := topoObj.GetAspect(adHocAspect); err != nil {
				log.Errorf("Failed getting aspect: %v", err)
				return nil, err
			}

			if srcPort, ok = adHocAspect.Properties["srcPort"]; ok {
				log.Infof("Source port is: %v", srcPort)
			}

			if dstPort, ok = adHocAspect.Properties["dstPort"]; ok {
				log.Infof("Destionation port is: %v", dstPort)
			}

			// Append src device with its port, and dst device with its port
			devicePortMap[srcDevice] = append(devicePortMap[srcDevice], srcPort)
			devicePortMap[dstDevice] = append(devicePortMap[dstDevice], dstPort)
		}
	}

	return devicePortMap, nil
}

// Cycle time extension is statically 0 now, not sure what to do about it yet
// Base time is statically 0 on both seconds and fractional seconds now, not sure what to do about it yet
// Takes in port and device IP, creates updates for admin-cycle-time-extension, admin-base-time, and config-change
func getFinalElems(port string, deviceIp string) []*pb.Update {
	// Build update for admin cycle time extension
	cycleTimeExtUpd := &pb.Update{
		Path: &pb.Path{
			Elem: []*pb.PathElem{
				{
					Name: "interfaces",
					Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
				},
				{
					Name: "interface",
					Key:  map[string]string{"name": port},
				},
				{
					Name: "gate-parameters",
					Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
				},
				{
					Name: "admin-cycle-time-extension",
					Key:  map[string]string{},
				},
			},
			Target: deviceIp,
		},
		Val: &pb.TypedValue{
			Value: &pb.TypedValue_UintVal{
				UintVal: 0, // This should maybe be calculated???
			},
		},
	}

	// Build update for admin base time seconds
	baseTimeSecondsUpd := &pb.Update{
		Path: &pb.Path{
			Elem: []*pb.PathElem{
				{
					Name: "interfaces",
					Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
				},
				{
					Name: "interface",
					Key:  map[string]string{"name": port},
				},
				{
					Name: "gate-parameters",
					Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
				},
				{
					Name: "admin-base-time",
					Key:  map[string]string{},
				},
				{
					Name: "seconds",
					Key:  map[string]string{},
				},
			},
			Target: deviceIp,
		},
		Val: &pb.TypedValue{
			Value: &pb.TypedValue_UintVal{
				UintVal: 0, // This should maybe be calculated???
			},
		},
	}

	// Build update for admin base time fractional seconds
	baseTimeFractionalSecondsUpd := &pb.Update{
		Path: &pb.Path{
			Elem: []*pb.PathElem{
				{
					Name: "interfaces",
					Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
				},
				{
					Name: "interface",
					Key:  map[string]string{"name": port},
				},
				{
					Name: "gate-parameters",
					Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
				},
				{
					Name: "admin-base-time",
					Key:  map[string]string{},
				},
				{
					Name: "fractional-seconds",
					Key:  map[string]string{},
				},
			},
			Target: deviceIp,
		},
		Val: &pb.TypedValue{
			Value: &pb.TypedValue_UintVal{
				UintVal: 0, // This should maybe be calculated???
			},
		},
	}

	// Build update for config change
	confChangeUpd := &pb.Update{
		Path: &pb.Path{
			Elem: []*pb.PathElem{
				{
					Name: "interfaces",
					Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
				},
				{
					Name: "interface",
					Key:  map[string]string{"name": port},
				},
				{
					Name: "gate-parameters",
					Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
				},
				{
					Name: "config-change",
					Key:  map[string]string{},
				},
			},
			Target: deviceIp,
		},
		Val: &pb.TypedValue{
			Value: &pb.TypedValue_BoolVal{
				BoolVal: true,
			},
		},
	}

	return []*pb.Update{cycleTimeExtUpd, baseTimeSecondsUpd, baseTimeFractionalSecondsUpd, confChangeUpd}
}

func getStatusChangeElems(port string, deviceIp string, numOfTrafficClassEntries int) []*pb.Update {
	// Build update for gate enabled
	gateEnabledUpd := &pb.Update{
		Path: &pb.Path{
			Elem: []*pb.PathElem{
				{
					Name: "interfaces",
					Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
				},
				{
					Name: "interface",
					Key:  map[string]string{"name": port},
				},
				{
					Name: "gate-parameters",
					Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
				},
				{
					Name: "gate-enabled",
					Key:  map[string]string{},
				},
			},
			Target: deviceIp,
		},
		Val: &pb.TypedValue{
			Value: &pb.TypedValue_BoolVal{
				BoolVal: true,
			},
		},
	}

	// Build update for admin gate states
	gateStatesUpd := &pb.Update{
		Path: &pb.Path{
			Elem: []*pb.PathElem{
				{
					Name: "interfaces",
					Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
				},
				{
					Name: "interface",
					Key:  map[string]string{"name": port},
				},
				{
					Name: "gate-parameters",
					Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
				},
				{
					Name: "admin-gate-states",
					Key:  map[string]string{},
				},
			},
			Target: deviceIp,
		},
		Val: &pb.TypedValue{
			Value: &pb.TypedValue_UintVal{
				UintVal: 255, // Statically set all gates to be open for their initial state
			},
		},
	}

	// Build update for admin control list length
	controlListLenUpd := &pb.Update{
		Path: &pb.Path{
			Elem: []*pb.PathElem{
				{
					Name: "interfaces",
					Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
				},
				{
					Name: "interface",
					Key:  map[string]string{"name": port},
				},
				{
					Name: "gate-parameters",
					Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
				},
				{
					Name: "admin-control-list-length",
					Key:  map[string]string{},
				},
			},
			Target: deviceIp,
		},
		Val: &pb.TypedValue{
			Value: &pb.TypedValue_UintVal{
				UintVal: uint64(numOfTrafficClassEntries),
			},
		},
	}

	return []*pb.Update{gateEnabledUpd, gateStatesUpd, controlListLenUpd}
}

func getGclElems(sched *schedule.Schedule, port string, deviceIp string) []*pb.Update {
	var updates []*pb.Update
	// For every traffic class, create an entry in the admin-control-list
	for index, trafficClass := range sched.TrafficClasses {
		// Build update for type of operation
		operationUpd := &pb.Update{
			Path: &pb.Path{
				Elem: []*pb.PathElem{
					{
						Name: "interfaces",
						Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
					},
					{
						Name: "interface",
						Key:  map[string]string{"name": port},
					},
					{
						Name: "gate-parameters",
						Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
					},
					{
						Name: "admin-control-list",
						Key:  map[string]string{"index": string(index)},
					},
					{
						Name: "operation-name",
						Key:  map[string]string{},
					},
				},
				Target: deviceIp,
			},
			Val: &pb.TypedValue{
				Value: &pb.TypedValue_StringVal{
					StringVal: "set-gates-states",
				},
			},
		}

		// Build update for gate states
		gateStateUpd := &pb.Update{
			Path: &pb.Path{
				Elem: []*pb.PathElem{
					{
						Name: "interfaces",
						Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
					},
					{
						Name: "interface",
						Key:  map[string]string{"name": port},
					},
					{
						Name: "gate-parameters",
						Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
					},
					{
						Name: "admin-control-list",
						Key:  map[string]string{"index": string(index)},
					},
					{
						Name: "sgs-params",
						Key:  map[string]string{},
					},
					{
						Name: "gate-states-value",
						Key:  map[string]string{},
					},
				},
				Target: deviceIp,
			},
			Val: &pb.TypedValue{
				Value: &pb.TypedValue_UintVal{
					UintVal: getGateStatesValue(trafficClass.Name),
				},
			},
		}

		// Build update for time interval
		timeIntervalUpd := &pb.Update{
			Path: &pb.Path{
				Elem: []*pb.PathElem{
					{
						Name: "interfaces",
						Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
					},
					{
						Name: "interface",
						Key:  map[string]string{"name": port},
					},
					{
						Name: "gate-parameters",
						Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
					},
					{
						Name: "admin-control-list",
						Key:  map[string]string{"index": string(index)},
					},
					{
						Name: "sgs-params",
						Key:  map[string]string{},
					},
					{
						Name: "time-interval-value",
						Key:  map[string]string{},
					},
				},
				Target: deviceIp,
			},
			Val: &pb.TypedValue{
				Value: &pb.TypedValue_UintVal{
					UintVal: getInterval(trafficClass.AssignedPortion, sched.GatingCycle),
				},
			},
		}

		updates = append(updates, operationUpd)
		updates = append(updates, gateStateUpd)
		updates = append(updates, timeIntervalUpd)
	}

	return updates
}

// Get gate state values if traffic class matches a predefined traffic class (best effort will open two gates)
func getGateStatesValue(trafficClassName string) uint64 {
	switch trafficClassName {
	case "isochronous":
		return 128
	case "cyclic-sync":
		return 64
	case "cyclic-async":
		return 32
	case "alarms-events":
		return 16
	case "config-diag":
		return 8
	case "network-control":
		return 4
	case "best-effort":
		return 3
	default:
		log.Errorf("Traffic class was not found: %v", errors.New("Traffic class did not match any of the predefined traffic classes, all gates will be closed!"))
	}

	return 0
}

func getInterval(assignedPercentage int32, gatingCycle float32) uint64 {
	return uint64(gatingCycle * 1000 * (float32(assignedPercentage) / 100))
}

// TODO: Make generic in the case of path for updates, it is currently statically made for our switches
// Takes in gating cycle in milliseconds and builds two updates for numerator and denominator
func getAdminCycleTimeElems(gatingCycle float32, port string, deviceIp string) []*pb.Update {
	// Create update element for numerator
	numeratorUpd := &pb.Update{
		Path: &pb.Path{
			Elem: []*pb.PathElem{
				{
					Name: "interfaces",
					Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
				},
				{
					Name: "interface",
					Key:  map[string]string{"name": port},
				},
				{
					Name: "gate-parameters",
					Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
				},
				{
					Name: "admin-cycle-time",
					Key:  map[string]string{},
				},
				{
					Name: "numerator",
					Key:  map[string]string{},
				},
			},
			Target: deviceIp,
		},
		Val: &pb.TypedValue{
			Value: &pb.TypedValue_IntVal{
				IntVal: int64(gatingCycle),
			},
		},
	}

	// Create update element for denominator
	denominatorUpd := &pb.Update{
		Path: &pb.Path{
			Elem: []*pb.PathElem{
				{
					Name: "interfaces",
					Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
				},
				{
					Name: "interface",
					Key:  map[string]string{"name": port},
				},
				{
					Name: "gate-parameters",
					Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
				},
				{
					Name: "admin-cycle-time",
					Key:  map[string]string{},
				},
				{
					Name: "denominator",
					Key:  map[string]string{},
				},
			},
			Target: deviceIp,
		},
		Val: &pb.TypedValue{
			Value: &pb.TypedValue_IntVal{
				IntVal: 1000, // 1000 ensures the gating cycle is in milliseconds
			},
		},
	}

	return []*pb.Update{numeratorUpd, denominatorUpd}
}
