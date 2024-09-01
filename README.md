# low-level-functions

Low-level functions for golang that help to avoid allocations

# remarks:

- MakeNoZero calls mallocgc with the false flag, which means that the allocated memory will not be initialized with zeros. This is the key point
  _this function will be useful where initializing the array with zeros is not important (for example, if you write data to the array immediately), your function can be much faster because it skips the initialization step_

- String:
  _In Go, strings are immutable, which means that once a string is created, its contents should not change. However, in this case, the string references the same byte array as the original slice_

- StringToBytes same concept as for String
