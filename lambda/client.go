package minecraft

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/pterm/pterm"
	"os"
	"phoenix/lambda/function"
	"phoenix/ligo"
	"phoenix/minecraft"
	"phoenix/minecraft/protocol"
	"phoenix/minecraft/protocol/packet"
	"time"
)

type Callback func(output *packet.CommandOutput) error
type Client struct {
	bot, operator string
	spaces        map[string]*function.Space
	vm            *ligo.VM
	conn          *minecraft.Conn
	config        PlotConfig
	callbacks     map[string]Callback
}

func (client *Client) StartConsole() {
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			if scanner.Scan() {

				line := scanner.Text()
				value, err := client.vm.Eval(line)
				if err != nil {
					pterm.Error.Println(err)
				} else {
					pterm.Info.Println("==> ", value.Value)
				}
			}
		}
	}()
}

func (client *Client) Init() {
	client.vm.Funcs["get"] = func(vm *ligo.VM, variable ...ligo.Variable) ligo.Variable {
		client.SendCommand("gamerule sendcommandfeedback true", func(output *packet.CommandOutput) error { return nil })
		err := client.SendCommand(fmt.Sprintf("execute %s ~ ~ ~ testforblock ~ ~ ~ air", client.operator), func(output *packet.CommandOutput) error {
			pos, _ := function.SliceAtoi(output.OutputMessages[0].Parameters)
			pterm.Info.Println(output.OutputMessages[0])
			if len(pos) != 3 {
				return errors.New("testforblock function have got wrong number of positions")
			} else {
				space := vm.Vars["space"].Value.(*function.Space)
				space.SetPointer(pos)
				pterm.Info.Println("Position got: ", pos)
				client.SendCommand("gamerule sendcommandfeedback false", func(output *packet.CommandOutput) error { return nil })
			}
			return nil
		})
		if err != nil {
			return vm.Throw(fmt.Sprintf("SendCommand: %s", err))
		} else {
			return ligo.Variable{Type: ligo.TypeNil, Value: nil}
		}
	}

	client.vm.Funcs["plot"] = func(vm *ligo.VM, variable ...ligo.Variable) ligo.Variable {
		workSpace := vm.Vars["space"].Value.(*function.Space)
		if variable[0].Type == ligo.TypeArray {
			vec := variable[0].Value.([]function.Vector)
			for _, v := range vec {
				client.config.block.name = vm.Vars["block"].Value.(string)
				client.config.block.data = byte(vm.Vars["data"].Value.(int64))
				err := client.SetBlock(function.AddVector(v, workSpace.GetPointer()))
				time.Sleep(time.Millisecond)
				if err != nil {
					return vm.Throw(fmt.Sprintf("setblock: Unable to setblock: %s", err))
				}
			}
		} else if variable[0].Type == ligo.TypeFloat {
			workSpace.Plot(variable[0].Value.(function.Vector))
		} else {
			return ligo.Variable{Type: ligo.TypeErr, Value: "plot function's first argument should be of a vector or vector slice type"}
		}
		return ligo.Variable{
			Type: ligo.TypeNil,
		}
	}
}

func (client *Client) SendCommand(command string, callback Callback) error {
	requestID := uuid.New()
	callbackID := uuid.New()
	commandRequest := &packet.CommandRequest{
		CommandOrigin: protocol.CommandOrigin{
			Origin:         protocol.CommandOriginPlayer,
			UUID:           callbackID,
			RequestID:      requestID.String(),
			PlayerUniqueID: 0,
		},
		CommandLine: command,
		Internal:    false,
	}
	client.callbacks[callbackID.String()] = callback
	return client.conn.WritePacket(commandRequest)
}

func (client *Client) SendCommandWO(command string) error {
	commandRequest := &packet.SettingsCommand{
		CommandLine:    command,
		SuppressOutput: false,
	}
	return client.conn.WritePacket(commandRequest)
}

func (client *Client) SendCommandNoCallback(command string) error {
	requestID := uuid.New()
	callbackID := uuid.New()
	commandRequest := &packet.CommandRequest{
		CommandOrigin: protocol.CommandOrigin{
			Origin:         protocol.CommandOriginPlayer,
			UUID:           callbackID,
			RequestID:      requestID.String(),
			PlayerUniqueID: 0,
		},
		CommandLine: command,
		Internal:    false,
	}
	return client.conn.WritePacket(commandRequest)
}

func (client *Client) Actionbar(target, text string) error {
	return client.SendCommandNoCallback(fmt.Sprintf("title %s actionbar %s", target, text))
}

func (client *Client) SetBlock(pos function.Vector) error {
	cmd := fmt.Sprintf("setblock %v %v %v %s %d", pos[0], pos[1], pos[2], client.config.block.name, client.config.block.data)
	return client.SendCommandNoCallback(cmd)
}

func (client *Client) Info(text ...string) error {
	return client.SendCommand(InfoRequest("@a", text...), func(output *packet.CommandOutput) error { return nil })
}

func (client *Client) Error(text ...string) error {
	return client.SendCommand(ErrorRequest("@a", text...), func(output *packet.CommandOutput) error { return nil })
}

func InfoRequest(target string, lines ...string) string {
	now := time.Now().Format("§6[15:04:05]§b INFO: ")
	var items []TellrawItem
	for _, text := range lines {
		msg := fmt.Sprintf("%v %v", now, text)
		items = append(items, TellrawItem{Text: msg})
	}
	final := &TellrawStruct{
		RawText: items,
	}
	content, _ := json.Marshal(final)
	cmd := fmt.Sprintf("tellraw %v %s", target, content)
	return cmd
}

func ErrorRequest(target string, lines ...string) string {
	now := time.Now().Format("§6[15:04:05]§c ERROR: ")
	var items []TellrawItem
	for _, text := range lines {
		msg := fmt.Sprintf("%v %v", now, text)
		items = append(items, TellrawItem{Text: msg})
	}
	final := &TellrawStruct{
		RawText: items,
	}
	content, _ := json.Marshal(final)
	cmd := fmt.Sprintf("tellraw %v %s", target, content)
	return cmd
}

type TellrawItem struct {
	Text string `json:"text"`
}

type TellrawStruct struct {
	RawText []TellrawItem `json:"rawtext"`
}
