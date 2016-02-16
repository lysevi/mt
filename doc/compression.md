based on http://www.vldb.org/pvldb/vol8/p1816-teller.pdf

# Compressing time stamps

1. The block header stores the starting time stamp, 
2. For subsequent time stamps, tn:
(a) Calculate the delta of delta: D = (tn - t(n-1)) 
(b) If D is zero, then store a single '0' bit
(c) If D is between [-63, 64], store '10' followed by
the value (7 bits)
(d) If D is between [-255, 256], store '110' followed by
the value (9 bits)
(e) if D is between [-2047, 2048], store '1110' followed
by the value (12 bits)
(f) Otherwise store '1111' followed by D using 32 bits

# Compressing values

We then encode these XOR�d values with the following
variable length encoding scheme:
1. The first value is stored with no compression         
2. If XOR with the previous is zero (same value), store  single '0' bit
3. When XOR is non-zero, calculate the number of leading and trailing zeros in the XOR, store bit �1' followed
by either a) or b):
(a) (Control bit '0') If the block of meaningful bits falls within the block of previous meaningful bits,
i.e., there are at least as many leading zeros and  as many trailing zeros as with the previous value,
use that information for the block position and just store the meaningful XORed value.
(b) (Control bit '1') Store the length of the number of leading zeros in the next 7 bits, then store the
length of the meaningful XORed value in the next 6 bits. Finally store the meaningful bits of the
XORed value.

# Flags compression
if flag equals previous value flag, then write '0' else write '1' and 63 bit of flag