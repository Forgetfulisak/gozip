package main

/*

           Extra               Extr               Extra
      Code Bits Length(s) Code Bits Lengths   Code Bits Length(s)
      ---- ---- ------     ---- ---- -------   ---- ---- -------
       257   0     3       267   1   15,16     277   4   67-82
       258   0     4       268   1   17,18     278   4   83-98
       259   0     5       269   2   19-22     279   4   99-114
       260   0     6       270   2   23-26     280   4  115-130
       261   0     7       271   2   27-30     281   5  131-162
       262   0     8       272   2   31-34     282   5  163-194
       263   0     9       273   3   35-42     283   5  195-226a
       264   0    10       274   3   43-50     284   5  227-257
       265   1  11,12      275   3   51-58     285   0    258
       266   1  13,14      276   3   59-66

   The extra bits should be interpreted as a machine integer
   stored with the most-significant bit first, e.g., bits 1110
   represent the value 14.

            Extra           Extra               Extra
       Code Bits Dist  Code Bits   Dist     Code Bits Distance
       ---- ---- ----  ---- ----  ------    ---- ---- --------
         0   0    1     10   4     33-48    20    9   1025-1536
         1   0    2     11   4     49-64    21    9   1537-2048
         2   0    3     12   5     65-96    22   10   2049-3072
         3   0    4     13   5     97-128   23   10   3073-4096
         4   1   5,6    14   6    129-192   24   11   4097-6144
         5   1   7,8    15   6    193-256   25   11   6145-8192
         6   2   9-12   16   7    257-384   26   12  8193-12288
         7   2  13-16   17   7    385-512   27   12 12289-16384
         8   3  17-24   18   8    513-768   28   13 16385-24576
         9   3  25-32   19   8   769-1024   29   13 24577-32768
*/

type Code struct {
	len      int
	bitCount int
}

var lengthCodes = map[int]Code{
	257: Code{bitCount: 0, len: 3},
	258: Code{bitCount: 0, len: 4},
	259: Code{bitCount: 0, len: 5},
	260: Code{bitCount: 0, len: 6},
	261: Code{bitCount: 0, len: 7},
	262: Code{bitCount: 0, len: 8},
	263: Code{bitCount: 0, len: 9},
	264: Code{bitCount: 0, len: 10},
	265: Code{bitCount: 1, len: 11},
	266: Code{bitCount: 1, len: 13},
	267: Code{bitCount: 1, len: 15},
	268: Code{bitCount: 1, len: 17},
	269: Code{bitCount: 2, len: 19},
	270: Code{bitCount: 2, len: 23},
	271: Code{bitCount: 2, len: 27},
	272: Code{bitCount: 2, len: 31},
	273: Code{bitCount: 3, len: 35},
	274: Code{bitCount: 3, len: 43},
	275: Code{bitCount: 3, len: 51},
	276: Code{bitCount: 3, len: 59},
	277: Code{bitCount: 4, len: 67},
	278: Code{bitCount: 4, len: 83},
	279: Code{bitCount: 4, len: 99},
	280: Code{bitCount: 4, len: 115},
	281: Code{bitCount: 5, len: 131},
	282: Code{bitCount: 5, len: 163},
	283: Code{bitCount: 5, len: 195},
	284: Code{bitCount: 5, len: 227},
	285: Code{bitCount: 0, len: 258},
}

var distanceCodes = map[int]Code{
	0:  Code{bitCount: 0, len: 1},
	1:  Code{bitCount: 0, len: 2},
	2:  Code{bitCount: 0, len: 3},
	3:  Code{bitCount: 0, len: 4},
	4:  Code{bitCount: 1, len: 5},
	5:  Code{bitCount: 1, len: 7},
	6:  Code{bitCount: 2, len: 9},
	7:  Code{bitCount: 2, len: 13},
	8:  Code{bitCount: 3, len: 17},
	9:  Code{bitCount: 3, len: 25},
	10: Code{bitCount: 4, len: 33},
	11: Code{bitCount: 4, len: 49},
	12: Code{bitCount: 5, len: 65},
	13: Code{bitCount: 5, len: 97},
	14: Code{bitCount: 6, len: 129},
	15: Code{bitCount: 6, len: 193},
	16: Code{bitCount: 7, len: 257},
	17: Code{bitCount: 7, len: 385},
	18: Code{bitCount: 8, len: 513},
	19: Code{bitCount: 8, len: 769},
	20: Code{bitCount: 9, len: 1025},
	21: Code{bitCount: 9, len: 1537},
	22: Code{bitCount: 10, len: 2049},
	23: Code{bitCount: 10, len: 3073},
	24: Code{bitCount: 11, len: 4097},
	25: Code{bitCount: 11, len: 6145},
	26: Code{bitCount: 12, len: 8193},
	27: Code{bitCount: 12, len: 12289},
	28: Code{bitCount: 13, len: 16385},
	29: Code{bitCount: 13, len: 24577},
}
