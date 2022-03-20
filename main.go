package main

import (
	"log"
	"os"

	"github.com/vapourismo/knx-go/knx"
	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/dpt"
	"github.com/vapourismo/knx-go/knx/util"
)

func main() {
	// Setup logger for auxiliary logging. This enables us to see log messages from internal
	// routines.
	util.Logger = log.New(os.Stdout, "", log.LstdFlags)

	// Connect to the gateway.
	client, err := knx.NewGroupTunnel("10.0.0.7:3671", knx.DefaultTunnelConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Close upon exiting. Even if the gateway closes the connection, we still have to clean up.
	defer client.Close()

	// Send 20.5°C to group 1/2/3.
	err = client.Send(knx.GroupEvent{
		Command:     knx.GroupWrite,
		Destination: cemi.NewGroupAddr3(1, 2, 3),
		Data:        dpt.DPT_9001(20.5).Pack(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Receive messages from the gateway. The inbound channel is closed with the connection.
	for msg := range client.Inbound() {
		var temp dpt.DPT_9001

		err := temp.Unpack(msg.Data)
		if err != nil {
			continue
		}

		util.Logger.Printf("%+v: %v", msg, temp)
	}
}
