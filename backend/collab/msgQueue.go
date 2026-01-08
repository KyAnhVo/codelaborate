package collab

type MsgQueue struct {
	// implementing circular msg array

	// list of UpdateMsg
	q []*UpdateMsg

	// left bound of q
	l int

	// right bound of q
	r int
}

func CreateMsgQueue() *MsgQueue {
	queue := MsgQueue{
		q: make([]*UpdateMsg, 32),
		l: 0,
		r: 0,
	}
	return &queue
}

func (queue *MsgQueue) Enqueue(msg *UpdateMsg) {
	// if at limit, double the size
	if queue.IsFull() {
		newq := make([]*UpdateMsg, cap(queue.q) * 2)
		r := 0
		for queue.l != queue.r {
			newq[r] = queue.q[queue.l]
			queue.l = (queue.l + 1) % cap(queue.q)
			r++
		}
		queue.q, queue.l, queue.r = newq, 0, r
	}

	// add item in
	queue.q[queue.r] = msg
	queue.r = (queue.r + 1) % cap(queue.q)
}

func (queue *MsgQueue) Dequeue() *UpdateMsg {
	if queue.IsEmpty() {
		return nil
	}

	msg := queue.q[queue.r]
	// if r != 0 then r := r - 1
	// if r == 0 then r := cap(q) - 1
	queue.r = queue.nextIndex(queue.r)
	return msg
}

func (queue *MsgQueue) IsEmpty() bool {
	return queue.l == queue.r
}

func (queue *MsgQueue) IsFull() bool {
	return (queue.r + 1) % cap(queue.q) == queue.l
}

func (queue *MsgQueue) nextIndex(ind int) int {
	queueCap := cap(queue.q)
	return (ind + queueCap - 1) % queueCap
}
