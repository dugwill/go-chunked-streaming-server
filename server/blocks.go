package server

import (
	"log"
	"time"
)

type cmafChunk struct {
	t    time.Time
	size int
}

type cteChunk struct {
	t    time.Time
	size int
}

func processBlocks(f *File) {

	var cteChunks []cteChunk
	//var cmafChunks []cmafChunk

	log.Printf("%s - Blocks Received: %d", f.Name, len(f.buffBlocks))

	// Apparently, on Windows, the max read buffer is 32768. Not sure if this is
	// Windows or GO. But the processing of the blocks is base on this Max number.
	// If the block is less then this
	// then that is a complete CTE Chunk.  If it is equal to that then the CTE Chunk size is
	// the sum of the size of this and the next block.  There could be a case where
	// the CTE Chunk is exactly this size and that is an issue right now.

	// Get CTE Chunks
	for i := 0; i < len(f.buffBlocks); i++ {
		if f.buffBlocks[i].size >= 32768 {
			cteChunks = append(cteChunks, cteChunk{f.buffBlocks[i].t, f.buffBlocks[i].size + f.buffBlocks[i+1].size})
			// increment count since we already add the next block size
			i++
		} else {
			cteChunks = append(cteChunks, cteChunk{f.buffBlocks[i].t, f.buffBlocks[i].size})
		}
	}
	log.Printf("%s - cteChunks: %d", f.Name, len(cteChunks))
	//for i := 0; i < len(cteChunks); i++ {
	//	log.Printf("%s - cteChunk: %d  Time %s  Size: %d", f.Name, i, cteChunks[i].t.String(), cteChunks[i].size)
	//}
}
