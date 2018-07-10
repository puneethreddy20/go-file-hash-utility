# Code Documentation
   
   The main goal is to build a command line utility to turn a UTF-8 input stream into a 64bit hash.
   The process reads an arbitrary long byte stream on the standard input and partition them into equally sized chunks.
   The chunk size is specified by the user and must be a multiple of 8.

   1. Procedure:
   
        ->> Step1: Window and Threads value
   
        While running the application, it takes two arguments window, threads. window value determines the size of
        each chunk and it must be multiple of 8. Default value of window is 32768.
   
        Threads value determines number of threads/go routines to be used. The default value is 1.
        And based on the input, we determine maximum possible threads value.
   
        For example: Numberofchunks=10, Threads value given by user=100
   
        In the above case, The program takes 10 as thread value. Because each thread/go routine processess atleast 
        1 chunk. As there are only 10 chunks, 10 threads would be enough.
   
        If Numberofchunks=100, User's threads value=50:
        In this case, The program takes 50 as thread value. some threads will processess multiple chunks(divided between threads), 
        but all threads will process atleast one chunk.
   
        The program first checks if the window and threads value are valid or not. If not valid it returns err to cmdline.
   
        Allchunks non-buffered channel is created.
        
        ->> Step2: Reading Input from stdin and processing dividing it into chunks based on window value.
        
        ioutil.ReadAll(os.Stdin): It allows us to get all the byte stream from stdin. All the input is got at once.
        This byte stream is converted into string.
   
        The String is divided into chunks; where each chunk has size of window's value. If the last chunk doesn't have
        enough bytes, its adds padding of byte value of 0. Thus, with padding the last chunk also has same size.
   
        All the chunks are stored in a string array. And is pushed into Allchunks non buffered channel.
        
        ->> Step3: Allchunks gets the value from ALLchunks-nonbuffered channel. and Create Go routines which process atleast one chunk.
   
        The main waits until Allchunks gets string array from Allchunks non buffered channel. 
   
        The following attributes are created or calculated:
   
        1. Number of chunks by finding length of Allchunks string array.
   
        2. Number of buffer channels (it is number of possible 64bit unsigned integers for allchunks with certain window size.)
   
        3. hashchannels is created (buffered channel of size as number of buffer channels.)
   
        4. Number of possible threads. It is minimum value of threads and number of chunks.
   
        5. An integer n, which stores the value of (NumofChunks/possibleThreads). It is used while slicing the chunks for a go routine input.
   
        6. Remainder of (NumofChunks/possibleThreads). It is used to distribute the chunks almost evenly to go routines.
   
        Using the above values go routines are created and 'hash' waits for go routines to push its outcome values 
        to hashchannels.
   
    for i := 0; i < possibleThreads; i++ {
   		r := x + n
   		if reminder>0 {
   			r++
   			reminder--
   		}
   		//Increment wg(wait group) before starting a go routines.
   		wg.Add(1)
   		go ChunkThread(Allchunks[x:r], windowsize, hashChannels, x, r,&wg)
   		x = r
   	}
   	
   
       
   
   The above code creates go routines which take Allchunks[x:r] (slice, which has multiple chunks), Allchunks is
   divided almost evenly for all threads.
       
   For example: 
   
   if possible threads=6335 but given threads is 6000.
   Then in that case: if we don't use remainder then 5999 threads would have slice with one chunk. 6000th thread
   would have 335 chunks to process.
   
   To avoid that we get the remainder of possible threads and given threads, for the above case the remainder 
   is 335. So for threads upto 335 would have 2 chunks and remaining will have 1 chunk. 
   

   ->> step4: The whole process of ChunkThread function. (logic happening in go routine.)

   The ChunkThread functions has Allchunks[startindex:endindex] slice, windowsize, hashChannels, startindex,endindex
   and address of the wg(wait group).
   
   For Each Chunk in Allchunks[startindex:endindex], 
   we call 'EndvalueofChunk' Function which has the parameters:
    particular chunk, its index (chunkindex)value according to original Allchunks(which can determined using startindex and endindex),windowsize and hashchannels.
   
   
   In EndvalueofChunk function and Convert 8 byte string to 64 bit unsigned integer (Convertchunckstrtouint64)function:
   
   EndvalueofChunk function returns the final value after hashing each chunk.
   The input string 's' which is a string chunk and its length is same as window size, 
   also in multiples of 8. Divide the input string 's' into multiple 8 byte strings.
  
  (Ex: if chunk length and window size is 16 then it can be divided into two 8 byte strings.)
   These 8 bytes strings are converted into 64 bit unsigned integers by Convertchunckstrtouint64 function.
   
   EndvalueChunk also has input called chunkindex which is the index of this chunk in  All CHUNKS array.
   Each integer in this chunck is logical or-ed with the chunk index and pushed into hashchannel buffer.
   
   why logical or-ed with chunk index?
   
   Because the Question said, each integer is logically or-ed with zero-based chunk index.
   
   Ex: chunkindex=2 chunk is converted to two 8 bytes strings. The two 8 byte strings are converted to two 64 bit integers
   Each 64-bit unsigned interger is logically or-ed with 2 (chunkindex) and pushed into hashchannel.
   
   
   
   ->> step5: MeanWhile What's happening in the main():

   EndvalueChunk also has input called chunkindex which is the index of this chunk in  All CHUNKS array.

       for i := 0; i < numberofBufferchannels; i++ {
            intermediateHash := <-hashChannels //wait for each chunk integers to be pushed into hashchannel buffer
            hash = hash ^ intermediateHash //xor-ed
        }
       
   
   As hash channel is buffered channel, because we know how many 64 bit unsigned integers are possible.
   (based on number of chunks and window size).
   
   
   It goes into for loop:
       
   1.waits for each chunk integers to be pushed into hashchannel buffer. and stored in intermediateHash
   
   2.hash is xor-ed with intermediateHash.
   
   we wait until all go routines are done.(using wg.Wait()) thus hash is printed into cmdline.

   As xoring is commutative,  every time  xor-ed with an 64 bit unsigned integer which appears randomly from any go routine, so each time we run it, we would get same hash value while using multi threading.