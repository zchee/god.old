// Copyright 2017 The god Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package god

import (
	serialpb "github.com/zchee/god/serial"
)

// CreateLocation create pb.Location.
func CreateLocation(filename string, offset int64) *serialpb.Location {
	return &serialpb.Location{
		Filename: filename,
		Offset:   offset,
	}
}
