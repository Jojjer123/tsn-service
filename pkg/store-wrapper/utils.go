package storewrapper

import (
	"context"
	"strings"
	"tsn-service/pkg/structures/configuration"

	"github.com/atomix/atomix-go-client/pkg/atomix"
	"github.com/golang/protobuf/proto"
)

// Takes in an object as a byte slice, a URN in the
// format of "storeName.Resource", and stores the
// structure at the URN.
func sendToStore(obj []byte, urn string) error {
	ctx := context.Background()

	// Create a slice of URN elements
	urnElems := strings.Split(urn, ".")

	// Get the store
	store, err := atomix.GetMap(ctx, urnElems[0])
	if err != nil {
		log.Errorf("Failed getting store \"%s\": %v", urnElems[0], err)
		return err
	}

	// TODO: Check if the URN contains more complex path and do something special then

	// Store the object
	_, err = store.Put(ctx, urnElems[1], obj)
	if err != nil {
		log.Errorf("Failed storing resource \"%s\": %v", urnElems[1], err)
		return err
	}

	return nil
}

func getFromStore(urn string) (*configuration.Request, error) {
	ctx := context.Background()

	// Create a slice of maximum two URN elements
	urnElems := strings.SplitN(urn, ".", 2)

	// log.Info("Getting map...")

	// Get the store
	store, err := atomix.GetMap(ctx, urnElems[0])
	if err != nil {
		log.Errorf("Failed getting store \"%s\": %v", urnElems[0], err)
		return &configuration.Request{}, err
	}

	// log.Info("Getting obj from store...")

	// TODO: Check if the URN contains more complex path and do something special then

	// Get the object from store
	obj, err := store.Get(ctx, urnElems[1])
	if err != nil {
		log.Errorf("Failed getting resource \"%s\": %v", urnElems[1], err)
		return &configuration.Request{}, err
	}

	// log.Info("Unmarshaling object...")

	// Unmarshal the byte slice from the store into request data
	var req = configuration.Request{}
	err = proto.Unmarshal(obj.Value, &req)
	if err != nil {
		log.Errorf("Failed to unmarshal request data from store: %v", err)
		return &configuration.Request{}, nil
	}

	return &req, nil
}
