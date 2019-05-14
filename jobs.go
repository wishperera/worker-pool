package worker_pool

import (
	"github.com/google/uuid"
	"context"
)

type Job struct {
	ctx context.Context
	id uuid.UUID
	input interface{}
	output interface{}
	err    error
}


/*
	Func NewJob(ctx context.Context,in interface{})(j *Job)
	------------------------------------------------------
	Returns a pointer to a Job object with provided input.
 */
func NewJob(ctx context.Context,in interface{})(j *Job){
	return &Job{
		id: uuid.New(),
		input: in,
		output:nil,
		err:nil,
	}
}


func (j *Job)GetOutput()(output interface{},err error){
	return  j.output,j.err
}

func (j *Job)GetInput()(input interface{}){
	return j.input
}

func (j *Job)GetID()(uuid.UUID){
	return j.id
}



