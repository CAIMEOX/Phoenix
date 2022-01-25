package minecraft

import (
	"errors"
	"fmt"
	"github.com/pterm/pterm"
	"os"
	"path/filepath"
	"phoenix/lambda/function"
	"phoenix/lambda/function/generator"
	"phoenix/lambda/function/std"
	"phoenix/ligo"
	"phoenix/minecraft"
	"phoenix/minecraft/auth"
	"phoenix/minecraft/protocol/packet"
)

type PlotConfig struct {
	block Block
}

type Block struct {
	name string
	data byte
}

func (client *Client) GetSpace(name string) *function.Space {
	return client.spaces[name]
}

func Run(path string) {
	config := function.ReadConfig(path)
	if config.Debug.Enabled {
		pterm.EnableDebugMessages()
	}

	// Init Connection :: Start
	dialer := func() minecraft.Dialer {
		if config.User.Auth {
			return minecraft.Dialer{
				TokenSource: auth.TokenSource,
			}
		} else {
			return minecraft.Dialer{}
		}
	}()

	conn, err := dialer.Dial("raknet", config.Connection.RemoteAddress)
	if err != nil {
		pterm.Error.Println(err)
	}
	defer conn.Close()
	// Init Connection :: End

	// Create a client by conn
	client := &Client{
		spaces:   make(map[string]*function.Space),
		vm:       ligo.NewVM(),
		bot:      config.User.Bot,
		operator: config.User.Operator,
		conn:     conn,
		config: PlotConfig{
			block: Block{
				name: "stone",
				data: 0,
			},
		},
		callbacks: make(map[string]Callback),
	}
	client.spaces["overworld"] = function.NewSpace()
	client.vm.Vars["space"] = ligo.Variable{
		Type:  ligo.TypeStruct,
		Value: client.spaces["overworld"],
	}

	if config.Lib.Std {
		std.StdInit(client.vm)
	}
	// Register the Generator Plugin
	generator.PluginInit(client.vm)
	if err := LoadScript(client.vm, config.Lib.Script); err != nil {
		pterm.Error.Println(err)
	}
	defaultConfig(client.vm)

	client.operator = config.User.Operator
	client.bot = config.User.Bot

	// Basic Functions Init
	client.Init()
	client.StartConsole()

	if err := conn.DoSpawn(); err == nil {
		pterm.Info.Println(fmt.Sprintf("Bot<%s> successfully spawned.", client.bot))
		// Collector : Get Position
		eval, err := client.vm.Eval(`(get)`)
		if err != nil {
			pterm.Error.Println(err)
		} else if eval.Value != nil {
			pterm.Info.Println(eval.Value)
		}
	} else {
		pterm.Error.Println(err)
	}

	// You will then want to start a for loop that reads packets from the connection until it is closed.
	for {
		// Read a packet from the connection: ReadPacket returns an error if the connection is closed or if
		// a read timeout is set. You will generally want to return or break if this happens.
		pk, err := conn.ReadPacket()
		if err != nil {
			break
		}

		// The pk variable is of type packet.Packet, which may be type asserted to gain access to the data
		// they hold:
		switch p := pk.(type) {
		case *packet.Text:
			if p.TextType == packet.TextTypeChat {
				if client.operator == p.SourceName {
					pterm.Info.Println(fmt.Sprintf("[%s] %s", p.SourceName, p.Message))
					value, err := client.vm.Eval(p.Message)
					if err != nil {
						pterm.Error.Println(err)
					} else {
						pterm.Info.Println("==> ", value.Value)
					}
				}
			}

		case *packet.CommandOutput:
			callback, ok := client.callbacks[p.CommandOrigin.UUID.String()]
			// TODO : Handle !ok
			if ok {
				delete(client.callbacks, p.CommandOrigin.UUID.String())
				if len(p.OutputMessages) > 0 {
					if !p.OutputMessages[0].Success {
						//pterm.Warning.Println(fmt.Sprintf("Unknown command: %s. Please check that the command exists and that you have permission to use it.", p.OutputMessages[0].Parameters))
					}
					err := callback(p)
					if err != nil {
						pterm.Warning.Println(err)
						// TODO : Handle error
					}
					continue
				}
			}
		case *packet.StructureTemplateDataResponse:
			data := p.StructureTemplate
			pterm.Info.Println(data)
		}

		// Write a packet to the connection: Similarly to ReadPacket, WritePacket will (only) return an error
		// if the connection is closed.
		p := &packet.RequestChunkRadius{ChunkRadius: 16}
		if err := conn.WritePacket(p); err != nil {
			break
		}

	}
}

func defaultConfig(vm *ligo.VM) {
	vm.Vars["block"] = ligo.Variable{
		Type:  ligo.TypeString,
		Value: "iron_block",
	}
	vm.Vars["data"] = ligo.Variable{
		Type:  ligo.TypeInt,
		Value: int64(0),
	}
}

func LoadScript(vm *ligo.VM, paths []string) error {
	for n, path := range paths {
		fileName := filepath.Base(paths[n])
		fileName = fileName[:len(fileName)-len(filepath.Ext(fileName))]
		if content, err := os.Open(path); err != nil {
			return err
		} else if err := vm.LoadReader(content); err != nil {
			return errors.New(fmt.Sprintf("Error loading script [%s]: %s", fileName, err))
		}
		pterm.Info.Println(fmt.Sprintf("Successfully loaded Script [%s] ", fileName))
	}
	return nil
}
