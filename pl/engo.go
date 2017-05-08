package pl

import (
	"log"
	"os"
	"reflect"
	"time"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

var vars *Vars
var scene, world Node

func (v *Vars) engo_bif_append() {
	var name string

	name = "engo-opts-type"
	v.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: engo_opts_type}))
	name = "engo-run"
	v.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: engo_run}))
}

func engo_run(v *Vars, args []Node) Node {
	vars = v
	scene = args[1]
	world = args[2]
	log.Println("engo_run started")
	engo.Run(args[0].(GoValueNode).v.Interface().(engo.RunOptions), &Scene{})
	log.Println("engo_run finished")
	v.raise_error("Ok")
	//vars.exit <- newIdentNode("Ok")
	log.Println("engo_run exit")
	time.Sleep(time.Millisecond * 100)
	os.Exit(0)
	return args[0]
}

func engo_opts_type(v *Vars, args []Node) Node {
	return newGoTypeNode(reflect.TypeOf(engo.RunOptions{}))
}

type Scene struct{}

type Entity struct {
	ecs.BasicEntity

	common.RenderComponent
	common.SpaceComponent
}

func (s *Scene) Preload() {
	//engo.Files.Load("icon.png")
	s.applyFunc("engo-scene-preload", []Node{scene})
}

func (s *Scene) Setup(w *ecs.World) {
	s.applyFunc("engo-scene-setup", []Node{scene, world})
	//common.SetBackground(color.White)

	//w.AddSystem(&common.RenderSystem{})

	// Retrieve a texture
	//texture, err := common.LoadedSprite("icon.png")
	//if err != nil {
	//	log.Println(err)
	//}

	// Create an entity
	//guy := Guy{BasicEntity: ecs.NewBasic()}

	// Initialize the components, set scale to 8x
	//guy.RenderComponent = common.RenderComponent{
	//	Drawable: texture,
	//	Scale:    engo.Point{8, 8},
	//}
	//guy.SpaceComponent = common.SpaceComponent{
	//	Position: engo.Point{0, 0},
	//	Width:    texture.Width() * guy.RenderComponent.Scale.X,
	//	Height:   texture.Height() * guy.RenderComponent.Scale.Y,
	//}

	// Add it to appropriate systems
	//for _, system := range w.Systems() {
	//	switch sys := system.(type) {
	//	case *common.RenderSystem:
	//		sys.Add(&guy.BasicEntity, &guy.RenderComponent, &guy.SpaceComponent)
	//	}
	//}
}

func (*Scene) Type() string { return "EngoScene" }

func (*Scene) applyFunc(name string, args []Node) {
	f := findFunc(newIdentNode(name), vars)
	if f != nil {
		applyFunc(f, args, vars)
	}
}
