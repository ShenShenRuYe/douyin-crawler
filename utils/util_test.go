/*
Package utils
@Time : 2023/1/9 20:06
@Author : 董胜烨
@File : util_test
@Software: GoLand
@note:
*/
package utils

import "testing"

func TestName(t *testing.T) {
	s := SaveFolder{}
	s.IsDateFolder = true
	s.Create()
}
