class PriorityQueue {
  constructor() {
    // min binary heap
    this.heap = [Math.MIN_SAFE_INTEGER];
    this.values = [undefined];
    this.length = 0;
  }

  internal_swap(a, b) {
    const priorityA = this.heap[a];
    const valueA = this.values[a];
    this.heap[a] = this.heap[b];
    this.values[a] = this.values[b];
    this.heap[b] = priorityA;
    this.values[b] = valueA;
  }

  // Recursive functions that reorder a branch
  internal_shakeout_down(index) {
    const leftChild = index * 2;
    const rightChild = leftChild + 1;
    const lesserChild =
      this.heap[leftChild] < this.heap[rightChild] ? leftChild : rightChild;
    if (this.heap[index] > this.heap[lesserChild]) {
      this.internal_swap(lesserChild, index);
      this.internal_shakeout_down(lesserChild);
    }
  }

  internal_shakeout_up(index) {
    const parent = Math.floor(index / 2);
    if (this.heap[parent] > this.heap[index]) {
      this.internal_swap(index, parent);
      this.internal_shakeout_up(parent);
    }
  }

  push(priority, item) {
    // // for O(n)
    // this.heap.push(item);
    // this.length++;
    // return;
    if (!Number.isInteger(priority)) {
      throw new Error("Priority must be integer");
    }
    let index = this.heap.length;
    this.heap.push(priority);
    this.values.push(item);
    this.internal_shakeout_up(index);
    this.length++;
  }

  pop() {
    // // for O(n)
    // this.length--;
    // return this.heap.pop();
    // return;
    if (this.heap.length == 1)
      throw new Error("Cannot pop, the queue is empty");
    const popped = this.values[1];
    this.internal_swap(1, this.heap.length - 1);
    this.heap.pop();
    this.values.pop();
    this.internal_shakeout_down(1);
    this.length--;
    return popped;
  }
}
