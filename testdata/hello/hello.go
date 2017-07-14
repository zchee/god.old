// Copyright 2017 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "fmt"

func main() {
	fn := fmt.Printf
	fn("hello")
	var fn2 = fmt.Printf
	fn2("hello2")
}
