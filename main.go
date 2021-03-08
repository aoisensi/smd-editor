package main

import (
	"math"
	"os"
	"reflect"

	"github.com/aoisensi/go-smd"
	cli "github.com/jawher/mow.cli"
)

func main() {
	app := cli.App("smd-editor", "edit smd file automaticaly")
	app.Command("add-anim-diff", "add animation diff file", addAnimDiff)
	app.Command("remove-shift", "remove shift", removeShift)
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func removeShift(cmd *cli.Cmd) {
	var (
		srcn = cmd.StringArg("SRC", "", "source file name")
		dstn = cmd.StringArg("DEST", "", "dest file name")
	)
	cmd.Action = func() {
		srcf, err := os.Open(*srcn)
		if err != nil {
			panic("src file not found")
		}
		defer srcf.Close()

		src, err := smd.Decode(srcf)
		if err != nil {
			panic(err)
		}

		max := len(src.Skeleton) - 1
		shift := sub3f(src.Skeleton[max][0].Pos, src.Skeleton[0][0].Pos)
		for i, bones := range src.Skeleton {
			bones[0].Pos = sub3f(bones[0].Pos, mul3f(shift, float64(i)/float64(max)))
		}

		dstf, err := os.Create(*dstn)
		if err != nil {
			panic(err)
		}
		defer dstf.Close()
		src.Encode(dstf)
	}
}

func addAnimDiff(cmd *cli.Cmd) {
	var (
		srcn  = cmd.StringArg("SRC", "", "source file name")
		diffn = cmd.StringArg("DIFF", "", "diff file name")
		dstn  = cmd.StringArg("DEST", "", "dest file name")
		//zero  = cmd.BoolOpt("zero", false, "zero frame only")
	)
	cmd.Action = func() {
		srcf, err := os.Open(*srcn)
		if err != nil {
			panic("src file not found")
		}
		defer srcf.Close()

		difff, err := os.Open(*diffn)
		if err != nil {
			panic("diff file not found")
		}
		defer difff.Close()

		src, err := smd.Decode(srcf)
		if err != nil {
			panic(err)
		}

		diff, err := smd.Decode(difff)
		if err != nil {
			panic(err)
		}

		if !reflect.DeepEqual(src.Nodes, diff.Nodes) {
			panic("skelton is not same")
		}

		dst := &smd.SMD{}
		dst.Nodes = src.Nodes
		dst.Skeleton = make(smd.Skeleton)
		frames := len(src.Skeleton)
		if len(diff.Skeleton) < frames {
			frames = len(diff.Skeleton)
		}
		for i := 0; i < frames; i++ {
			bones := make([]*smd.SkeletonBone, 0, 100)
			sbone := src.Skeleton[i][0]
			for _, dbone := range diff.Skeleton[i] {
				if sbone.BoneID != dbone.BoneID {
					continue
				}
				drot := dbone.Rot
				for _, node := range diff.Nodes {
					if node.ID == dbone.BoneID && node.ParentID == -1 {
						drot = [3]float64{drot[0], drot[1], drot[2] + math.Pi/2}
						break
					}
				}
				bone := &smd.SkeletonBone{
					BoneID: sbone.BoneID,
					Pos:    add3f(sbone.Pos, dbone.Pos),
					Rot:    combineEuler(sbone.Rot, drot),
				}
				bones = append(bones, bone)
				break
			}
			dst.Skeleton[i] = bones
		}

		dstf, err := os.Create(*dstn)
		if err != nil {
			panic(err)
		}
		defer dstf.Close()
		dst.Encode(dstf)
	}
}
