package fileiter

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type Iterator struct {
	head         *node
	tail         *node
	breadthFirst bool
}

type node struct {
	dir  string
	stat os.FileInfo
	err  error
	skip bool
	next *node
}

type IterationOption = func(*Iterator)

func BreadthFirst() IterationOption {
	return func(w *Iterator) {
		w.breadthFirst = true
	}
}

func New(dir string, opts ...IterationOption) *Iterator {
	stat, err := os.Lstat(dir)
	n := &node{dir: filepath.Dir(dir), stat: stat, err: err}
	w := &Iterator{head: n, tail: n}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func (w *Iterator) Next() bool {
	if w.head == nil {
		return false
	}
	if !w.head.skip && w.head.err == nil && w.head.stat.IsDir() {
		path := filepath.Join(w.head.dir, w.head.stat.Name())
		if ff, err := ioutil.ReadDir(path); len(ff) != 0 { // sorts ASC
			if err == nil {
				if w.breadthFirst {
					// add one-by-one to the tail
					for _, f := range ff {
						w.tail.next = &node{path, f, nil, false, nil}
						w.tail = w.tail.next
					}
					return w.next(w.head.next)
				} else {
					w.head = w.head.next // drop current
					// add one-by-one to the head (in reverse order)
					for i := len(ff) - 1; i > -1; i-- {
						w.head = &node{path, ff[i], nil, false, w.head}
					}
					return w.next(w.head)
				}
			} else {
				if w.breadthFirst {
					w.tail.next = &node{w.head.dir, w.head.stat, err, false, nil}
					w.tail = w.tail.next
					w.head = w.head.next // drop current
					return true
				} else {
					w.head.err = err
					return true
				}
			}
		}
	}
	return w.next(w.head.next)
}

func (w *Iterator) next(v *node) bool {
	w.head = v
	if v == nil {
		w.tail = nil
		return false
	}
	return true
}

func (w *Iterator) Dir() string {
	return w.head.dir
}

func (w *Iterator) Name() string {
	return w.head.stat.Name()
}

func (w *Iterator) IsDir() bool {
	return w.head.stat.IsDir()
}

func (w *Iterator) Err() error {
	return w.head.err
}

func (w *Iterator) SkipDir() {
	w.head.skip = true
}
