built a CRUD server based on this youtube video: https://www.youtube.com/watch?v=5BIylxkudaE

it did not include an update endpoint, so i decided to add one myself.

this line of the function was a little confusing to me because of the unfamiliar syntax

```go
s.shoppingItems = append(s.shoppingItems[:i], s.shoppingItems[i+1:]...)
```

i asked chat gpt to summarize it in detail, and this is what it came up with: 

# Chat GPT description of slices and `append`

Let's dive deeper into the line that actually removes an item from the slice:

```go
s.shoppingItems = append(s.shoppingItems[:i], s.shoppingItems[i+1:]...)
```
This line uses two main concepts: slicing and appending in Go. Both of these operations are crucial for understanding how the item is removed from the slice.

1. Slicing in Go:
In Go, slices are like dynamic arrays, and you can create sub-slices from an existing slice using the slice operator [:]. The slice operator allows you to create a new slice that refers to a portion of an existing slice.

The syntax is:

```go
slice[start:end]
```

Where:
start is the index where the slice starts (inclusive).
end is the index where the slice ends (exclusive).
If start is omitted, it defaults to 0, and if end is omitted, it defaults to the length of the slice.

Let's break down how this works in your specific case:
```go
s.shoppingItems[:i]:
```
This takes all elements from the start of the shoppingItems slice up to, but not including, index i.
In Go, when you slice [:i], the resulting slice includes elements from index 0 to i-1.
This is the part of the slice before the item you want to remove.
```go
s.shoppingItems[i+1:]:
```
This takes all elements from index i+1 to the end of the shoppingItems slice.
In Go, when you slice [i+1:], the resulting slice includes elements starting from index i+1 to the last element of the slice.
This is the part of the slice after the item you want to remove.
The item at index i (the one you want to remove) is excluded because it is neither included in [:i] nor in [i+1:].

Example:
If you have the following slice:

```go
s.shoppingItems = []Item{
    {ID: 1, Name: "Apple"},
    {ID: 2, Name: "Banana"},
    {ID: 3, Name: "Orange"},
}
```
And you want to remove the item at index i = 1 (the "Banana"):

```go
s.shoppingItems[:i] gives:
```
```go
[]Item{
    {ID: 1, Name: "Apple"},
}
```
This includes everything before "Banana".

```go
s.shoppingItems[i+1:] gives:
```
```go
[]Item{
    {ID: 3, Name: "Orange"},
}
```
This includes everything after "Banana".

1. Appending in Go:
The append function in Go creates a new slice by adding elements to the end of an existing slice. It can take multiple elements or another slice as arguments.

The syntax is:

```go
append(existingSlice, newElements...)
```
existingSlice is the slice to which new elements will be appended.
newElements can be individual elements or another slice.
The ... operator is used when appending a slice to another slice, which "unpacks" the elements of the slice you're appending.
Example:
Suppose you have:

```go
s.shoppingItems[:i]  // slice before "Banana"
[]Item{
    {ID: 1, Name: "Apple"},
}
```
And:

```go
s.shoppingItems[i+1:]  // slice after "Banana"
[]Item{
    {ID: 3, Name: "Orange"},
}
```
You can combine these two slices using append:

```go
s.shoppingItems = append(s.shoppingItems[:i], s.shoppingItems[i+1:]...)
```
This will append the second slice (s.shoppingItems[i+1:]) to the first slice (s.shoppingItems[:i]), resulting in:

```go
s.shoppingItems = []Item{
    {ID: 1, Name: "Apple"},
    {ID: 3, Name: "Orange"},
}
```
The item at index i ("Banana") has effectively been removed from the slice.

3. Complete Removal Process:
In summary, when the line s.shoppingItems = append(s.shoppingItems[:i], s.shoppingItems[i+1:]...) is executed:

Step 1: s.shoppingItems[:i] extracts all elements before the one at index i.
Step 2: s.shoppingItems[i+1:] extracts all elements after the one at index i.
Step 3: append combines these two sub-slices into a new slice, excluding the element at index i.
Step 4: The new slice (without the element at index i) is assigned back to s.shoppingItems, effectively removing the item.
Key Points to Note:
Memory Efficiency: Slices in Go are backed by arrays, so when you slice a slice, the underlying array is shared. The original slice is not immediately copied, and memory for the removed element isn't freed until the slice's capacity is reduced (which happens if the capacity changes or the slice is resliced later).
Immutability: In Go, slices are not immutable, but this operation creates a new slice with the element removed, leaving the original memory untouched.
Order of Elements: The order of elements in the slice is preserved except for the element at index i, which is removed.
Example Visualization:
Before:

```go
s.shoppingItems = [{ID: 1, Name: "Apple"}, {ID: 2, Name: "Banana"}, {ID: 3, Name: "Orange"}]
```
After removing the item at index 1 ("Banana"):

```go
// Index:        0          1
s.shoppingItems = [{ID: 1, Name: "Apple"}, {ID: 3, Name: "Orange"}]
```
Recap:
s.shoppingItems[:i] gives you the slice up to the element you want to remove (exclusive).
s.shoppingItems[i+1:] gives you the slice starting right after the element you want to remove.
append(s.shoppingItems[:i], s.shoppingItems[i+1:]...) joins these two slices together, excluding the item at index i.
The result is a new slice with the element removed, which is then reassigned to s.shoppingItems.