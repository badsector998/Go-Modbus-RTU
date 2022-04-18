package main

import (
	"fmt"
	"time"

	ser "github.com/goburrow/serial"
	mod "github.com/things-go/go-modbus"
)

func main() {
	fmt.Println("Modbus Poll Raw")

	var port, parity string
	var baudRate, dataBits, stopBits int
	var err error
	var res []byte

	port = "COM1"
	parity = "N"
	baudRate = 9600
	dataBits = 8
	stopBits = 1

	newClient := createNewHandler(port, baudRate, dataBits, stopBits, parity)
	err = newClient.Connect()
	defer newClient.Close()

	if err != nil {
		fmt.Println("connection refused", err)
		delay(10)
		return
	}

	var slaveId byte
	var address, quantity uint16

	slaveId = 8
	address = 10
	quantity = 10

	res, err = newClient.ReadHoldingRegistersBytes(slaveId, address, quantity)

	if err != nil {
		fmt.Println("Measuring error", err)
		delay(10)
		return
	}

	var result uint64
	for i := 0; i < int(2*quantity); i++ {
		result = result<<8 + uint64(res[i])
	}

	fmt.Println(result)
}

func delay(dur int) {
	duration := time.Duration(dur) * time.Second
	time.Sleep(duration)
}

func createNewHandler(port string, baudRate, dataBits, stopBits int, parity string) mod.Client {
	p := mod.NewRTUClientProvider(
		mod.WithSerialConfig(ser.Config{
			Address:  port,
			BaudRate: baudRate,
			DataBits: dataBits,
			StopBits: stopBits,
			Parity:   parity,
			Timeout:  10 * time.Second,
		}),
	)

	client := mod.NewClient(p)
	return client
}
