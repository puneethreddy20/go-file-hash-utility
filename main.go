package main

import (
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sync"
)

/*
This Block allows us to takes cmdline arguments for window and threads.
*/
var (
	threads = flag.Int("threads", 1, "Number of threads")

	window = flag.Int("window", 32768, "window size")
)

/*
we will see if the given cmdline argument window's value is valid or not.
It should be multiple of 8 and greater than 0.
*/

func IsWindowValid(windowsize int) (bool, error) {
	if windowsize > 0 && windowsize%8 == 0 {
		return true, nil
	}
	log.Println("Window size should be greater than zero and in multiples of 8")
	return false, errors.New("Window size should be greater than zero and in multiples of 8")
}

/* This block checks if the given thread value is valid or not.
thread value should be greater than 0.
*/
func IsThreadsnumValid(threads int) (bool, error) {
	if threads <= 0 {
		log.Println("Threads should be greater than 0")
		return false, errors.New("Threads should be greater than 0")
	}
	return true, nil
}

/*
This function adds padding bytes of 0 to the chunk
which has length less than window size
(i.e. if chunk length<window size then add padding)
*/
func AddPadding(chunckstring string, size int) (string, error) {
	chunckbytes := []byte(chunckstring)
	length := len(chunckbytes)
	if length == size {
		return chunckstring, nil
	}
	if length > size {
		return "", errors.New("Wrong call (add padding is supposed to be called only when chunckbytes " +
			"size is less than or equal to the given size)")
	}
	diff := size - length
	for i := 0; i < diff; i++ {
		chunckbytes = append(chunckbytes, byte(0))
	}
	return string(chunckbytes), nil
}

/*
This functions returns all the chunks of in an string array,
each chunk length is equal to the given window size. it also adds padding to string, if its length is below window size.
But that probably happens only for the last string.
*/
func GetAllchuncks(s string, windowsize int) ([]string, error) {
	var chunks []string
	runes := []rune(s) //type rune
	runeslength := len(runes)
	if runeslength == 0 {
		return []string{s}, nil
	}
	for i := 0; i < runeslength; i = i + windowsize {
		n := i + windowsize
		if n > runeslength {
			paddedstring, err := AddPadding(string(runes[i:runeslength]), windowsize)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			chunks = append(chunks, paddedstring)
			continue
		}
		chunks = append(chunks, string(runes[i:n]))
	}
	return chunks, nil
}

/*
This functions converts chunk string to unsigned int 64 bits. Each string is probably 8 bytes.
*/
func Convertchunckstrtouint64(s string) (uint64, error) {
	chunckbytes := []byte(s)
	//chunckbytes1:=[]byte{0,0,0,0,0,0,0,66}
	var a uint64
	for i := 0; i < len(chunckbytes); i++ {
		slice := chunckbytes[i : i+1]
		//fmt.Println(slice)
		b, _ := binary.Uvarint(slice)
		a = a*256 + b
	}
	return a, nil
}

/*
This function returns the final value after hashing each chunk. The input string 's' which is a
string chunk and its length is same as window size, also in multiples of 8. Divide the input string 's' into
multiple 8 byte strings.(Ex: if chunk length and window size is 16 then it can be divided into two 8 byte strings.
These 8 bytes strings are converted into 64 bit unsigned integers.

This functions also has input called chunkindex which is the index of this chunk in  All CHUNKS array.
Each integer in this chunck is logical or-ed with the chunk index and pushed into hashchannel buffer.
*/
func EndvalueofChunk(s string, windowsize int, hashchannel chan uint64, chunckindex int) (uint64, error) {

	length := len(s)
	if length != windowsize {
		log.Printf("chunk string length and window size are not eqaul! or may be its an empty file")
		return 0, errors.New("chunk string length and window size are not eqaul")
	}
	var finalvalue uint64
	var value uint64
	var err error
	val := uint64(chunckindex)

	for i := 0; i < length; i += 8 {
		value, err = Convertchunckstrtouint64(s[i : i+8])
		if err != nil {
			log.Println(err)
			return 0, err
		}
		finalvalue = value | val
		hashchannel <- finalvalue
	}
	return finalvalue, nil
}

/*
So based on number of threads given by the end user, ALL CHUNKS array is sliced and inputted to go routines
which are run concurrently.

This function gets the sliced array, windowsize, starting index, end index, hashchannel, wait group.
For each chunk string in sliced array, a function EndvalueofChunk which takes chunk string,
hashchannel,windowsize and chunk index(in ALL CHUNKS array, which can be determined by starting index value)
is called. Once all strings in slice are processed we end with wg.Done() which indicates that go routine is done.


*/
func ChunkThread(chunk []string, windowsize int, hashchannel chan uint64, index int, endindex int, wg *sync.WaitGroup) error {
	for _, eachvalue := range chunk {
		if index > endindex {
			log.Println("oops this is not supposed to happen")
			return errors.New("oops this is not supposed to happen")
		}
		_, err := EndvalueofChunk(eachvalue, windowsize, hashchannel, index)
		if err != nil {
			log.Println(err)
			return err
		}
		index++
	}
	//decrement by 1.
	defer wg.Done()
	return nil
}

/*
A small function to find minimum value.
*/
func minimum(val1 int, val2 int) int {
	if val1 > val2 {
		return val2
	}
	return val1
}

func main() {
	flag.Parse()

	var wg sync.WaitGroup

	//check if window size is valid or not
	isWindowsizevalid, err := IsWindowValid(*window)
	if err != nil || !isWindowsizevalid {
		return
	}

	//check if number of threads is valid or not.
	isThreadsValid, err := IsThreadsnumValid(*threads)
	if err != nil || !isThreadsValid {
		return
	}

	Numofthreads := *threads
	windowsize := *window

	//creating non buffered channel which gets ALL CHUNKS string array, once processed by GetAllchuncks function.
	AllChunksChannel := make(chan []string)

	go func() {
		input, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
			return
		}

		chunks, err := GetAllchuncks(string(input), *window)
		if err != nil {
			log.Println("Error in GetAllchunks:", err)
			return
		}
		//send chunks to AllChunks channel
		AllChunksChannel <- chunks

	}()

	//initially hash is zero.
	var hash uint64 = 0

	//receive allchunks from channel.
	Allchunks := <-AllChunksChannel

	NumofChunks := len(Allchunks)

	//it is number of possible integers that are xor-ed with hash.
	numberofBufferchannels := ((windowsize) / 8) * NumofChunks

	hashChannels := make(chan uint64, numberofBufferchannels)

	possibleThreads := minimum(NumofChunks, Numofthreads)

	//n & remainder are factors used to slice the Allchunks array almost evenly between threads/go routines.
	n := NumofChunks / possibleThreads
	reminder := NumofChunks % possibleThreads

	x := 0

	for i := 0; i < possibleThreads; i++ {
		r := x + n
		if reminder > 0 {
			r++
			reminder--
		}
		//Increment wg(waitgroup) before starting a go routines.
		wg.Add(1)
		//run as go routine.
		go ChunkThread(Allchunks[x:r], windowsize, hashChannels, x, r, &wg)
		x = r
	}

	for i := 0; i < numberofBufferchannels; i++ {

		//wait for each chunk integers to be pushed into hashchannel buffer
		intermediateHash := <-hashChannels

		hash = hash ^ intermediateHash //xor-ed
	}

	//wg.Add increments count by 1 and wg.Done decrements count by 1. it waits until count is zero
	//which means until all go routines are done.
	wg.Wait()

	fmt.Printf("0x%X\n", hash)
	runtime.Gosched()
	close(AllChunksChannel)
	close(hashChannels)
}
