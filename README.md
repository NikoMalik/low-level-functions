# low-level-functions

> Low-level functions for golang that help to avoid allocations

1. **String(b []byte) string**

   - **What it does**: Converts a slice of bytes to a string without copying data, using insecure data access.
   - **Why it's faster**: Avoids copying data by creating a string that references the same bytes as the original slice.
   - **Risks**: Violates the immutability of strings in Go. If the source slice is changed, the string will change as well, which can lead to unexpected bugs.

2. **StringToBytes(s string) []byte**

   - **What it does**: Converts a string to a slice of bytes without copying the data, using insecure data access.
   - **Why it's faster**: Similar to `String`, avoids copying data by returning a slice that references the same bytes as the string.
   - **Risks**: If the string changes, the byte slice will change synchronously, which can cause problems when the data changes.

3. **CopyString(s string) string**

   - **What it does**: Creates a copy of a string by creating a new byte slice and copying the data from the string into that slice.
   - **Why it's faster**: Uses `MakeNoZero`, which does not initialize memory with zeros, speeding up the creation of a new slice.
   - **Risks**: Fewer risks compared to `String`, but there may be problems if the original data in the string must be kept intact.

4. **ConvertSlice[TFrom, TTo any](from []TFrom) ([]TTo, error)**

   - **What it does**: Converts a slice of one type to a slice of another type, provided the element sizes are the same.
   - **Why it's faster**: Doesn't create a new slice, just changes the slice header, keeping a reference to the same data in memory.
   - **Risks**: Conversion errors can cause failures if element sizes do not match or types are not compatible.

5. **MakeNoZero(l int) []byte**

   - **What it does**: Allocates memory for a byte slice without initializing it with zeros.
   - **Why it's faster**: Skips the step of initializing memory with zeros, which significantly speeds up the memory allocation process.
   - **Risks**: May lead to unexpected results if not all bytes are written before they are read.

6. **MakeNoZeroCap(l int, c int) []byte**

   - **What it does**: Creates a slice of bytes with a specific length and capacity without initializing with zeros.
   - **Why it's faster**: Similar to `MakeNoZero`, avoids initializing memory with zeros, which increases the speed of memory allocation.
   - **Risks**: Similar to `MakeNoZero`.

7. **StringBuffer (structure and methods)**

   - **What it does**: Implements a buffer for working with strings with the ability to change the length and capacity.
   - **Why it's faster**: Uses `MakeNoZeroCap` to allocate memory, which speeds up buffer handling. Also uses low-level methods to add data.
   - **Risks**: The main risks are associated with the use of `MakeNoZeroCap`; data initialization may be skipped.

8. **ConvertOne[TFrom, TTo any](from TFrom) (TTo, error)**

   - **What it does**: Converts a value of one type to a value of another type if their sizes are the same.
   - **Why it's faster**: Converts data directly through unsafe type conversion without creating additional data.
   - **Risks**: Conversion errors can result in incorrect data or crashes.

9. **MustConvertOne[TFrom, TTo any](from TFrom) TTo**

   - **What it does**: Converts one value of one type to a value of another type if their sizes match. Does not return an error.
   - **Why it's faster**: Avoids error checking, which makes the function execution faster.
   - **Risks**: If errors occur, there is no way to handle them, which may cause the program to crash.

10. **MakeZero[T any](l int) []T**

    - **What it does**: Allocates memory for a slice of any type with initialization with zeros.
    - **Why it's faster**: Optimized for large amounts of data, uses `mallocgc` with initialization.
    - **Risks**: Similar `MakeNoZero`.

11. **MakeZeroCap[T any](l int, c int) []T**

    - **What it does**: Creates a slice with a specific length and capacity with initialization with zeros.
    - **Why it's faster**: Similar to `MakeZero`, optimized for big data.
    - **Risks**: Similar `MakeNoZeroCap`.

12. **MakeNoZeroString(l int) []string**

    - **What it does**: Allocates memory for a slice of strings without initializing it with zeros.
    - **Why it's faster**: Skips initialization of memory with zeros.
    - **Risks**: Similar to `MakeNoZero`, may lead to unexpected results if strings are read before writing.

13. **MakeNoZeroCapString(l int, c int) []string**

    - **What it does**: Creates a slice of strings with a specific length and capacity without initializing with zeros.
    - **Why it's faster**: Skips memory initialization with zeros.
    - **Risks**: Similar to `MakeNoZeroCap`.

14. **Equal(a, b []byte) bool**

    - **What it does**: Compares two byte slices for equality using the low-level `memequal` function.
    - **Why it's faster**: Uses direct memory comparison, which is faster than byte-by-byte comparison.
    - **Risks**: If the slice lengths do not match, it immediately returns false, which may not always be the desired behavior.

15. **IsNil(v any) bool**

    - **What it does**: Checks if the value is nil.
    - **Why it's faster**: Uses direct checking with `reflect.ValueOf`.
    - **Risks**: Can cause panic when passing an uninitialized interface.

16. **IsEqual(v1, v2 any) bool**

    - **What it does**: Checks if two values point to the same memory location.
    - **Why it's faster**: Uses direct pointer comparison.
    - **Risks**: May not work correctly if values are of composite types or if objects are in different memory locations but with the same values.

17. **AtomicCounter (structure and methods)**

    - **What it does**: Implements an atomic counter using cache-line alignment to prevent false splitting.
    - **Why it's faster**: Avoids false splits by using cache-line alignment, which improves performance in multi-threaded environments.
    - **Risks**: More memory consumption due to alignment.

18. **GetItem[T any](slice []T, idx int) T**

    - **What it does**: Gets the slice item by index using an unsafe pointer cast.
    - **Why it's faster**: Avoids additional checks and accesses memory directly.
    - **Risks**: May cause panic if the index goes outside the slice.

19. **GetItemWithoutCheck[T any](slice []T, idx int) T**

    - **What it does**: Similar to `GetItem`, but without the index check.
    - **Why it's faster**: Removes the index check, which speeds up access.
    - **Risks**: Very unsafe, as going beyond the slice can cause data corruption or program crash.

20. **Pointer[T any](d T) \*T**
    - **What it does**: Returns a pointer to the passed value.
    - **Why it's faster**: Very simple operation, no data copying required.
    - **Risks**: Few risks, but requires caution when working with pointers.
