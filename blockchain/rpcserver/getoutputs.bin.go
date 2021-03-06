// Copyright 2017-2018 DERO Project. All rights reserved.
// Use of this source code in any form is governed by RESEARCH license.
// license can be found in the LICENSE file.
// GPG: 0F39 E425 8C65 3947 702A  8234 08B2 0360 A03A 9DE8
//
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY
// EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
// PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
// STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF
// THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package rpcserver

import "strconv"
import "net/http"

import "compress/gzip"

//import "github.com/pierrec/lz4"

// serve the outputs in streaming mode

// feeds any outputs requested by the server
func getoutputs(rw http.ResponseWriter, req *http.Request) {
	var err error
	start := uint64(0)
	stop := uint64(0)

	{ // parse start query parameter
		keys, ok := req.URL.Query()["start"]
		if !ok || len(keys) < 1 {
			//log.Println("Url Param 'key' is missing")
			//return
		} else {
			start_string := keys[0]
			start, err = strconv.ParseUint(start_string, 10, 64)
			if err != nil {
				start = 0
			}
		}
	}

	{ // parse stop query parameter
		keys, ok := req.URL.Query()["stop"]
		if !ok || len(keys) < 1 {

		} else {
			stop_string := keys[0]
			stop, err = strconv.ParseUint(stop_string, 10, 64)
			if err != nil {
				stop = 0
			}
		}
	}

	// do sanity check of stop  first
	top_id := chain.Get_Top_ID()
	biggest_output_index := chain.Block_Count_Vout(top_id) + chain.Get_Block_Output_Index(top_id)

	if stop == 0 || stop > biggest_output_index {
		stop = biggest_output_index
	}

	// feed in atleast 1 index
	if start >= stop {
		start = stop - 1
	}

	/*   lz4writer := lz4.NewWriter(rw)
	lz4writer.HighCompression = true // enable extreme but slow compression
	lz4writer.BlockMaxSize = 256*1024 // small block size to decrease memory consumption
	*/
	gzipwriter := gzip.NewWriter(rw)
	defer gzipwriter.Close()
	for i := start; i <= stop; i++ {
		// load the bytes and send them
		data, err := chain.Read_output_index(i)
		if err != nil {
			logger.Warnf("err while reading output err: %s\n", err)
			break
		}

		//
		//rw.Write(data)
		// lz4writer.Write(data)
		gzipwriter.Write(data)

	}
	//lz4writer.Flush() // flush any pending data
}
