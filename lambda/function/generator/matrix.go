package generator

import (
	"gonum.org/v1/gonum/mat"
	"phoenix/lambda/function"
	"phoenix/ligo"
)

func Dot(vm *ligo.VM, a ...ligo.Variable) ligo.Variable {
	if a[0].Type == ligo.TypeArray && a[1].Type == ligo.TypeArray {
		vec1 := a[0].Value.([]function.Vector)
		vec2 := a[1].Value.([]function.Vector)
		vecA := mat.NewDense(3, 3, Union(vec1))
		var r, vecB *mat.Dense
		if len(vec2[0]) == 2 {
			vecB = mat.NewDense(2, 2, Union(vec1))
		} else {
			vecB = mat.NewDense(3, 3, Union(vec2))
		}
		r.Mul(vecA, vecB)
		return vm.Throw("undone")
	} else {
		return vm.Throw("dot: first argument must be a Matrix")
	}
}

func Union(v []function.Vector) []float64 {
	var res []float64
	for _, vv := range v {
		res = append(res, vv...)
	}
	return res
}

func Pack(f []float64) []function.Vector {
	var res []function.Vector
	for i := 0; i < len(f); i += 3 {
		res = append(res, []float64{f[i], f[i+1], f[i+2]})
	}
	return res
}
