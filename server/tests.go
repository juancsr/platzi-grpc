package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/juancsr/platzi-grpc/models"
	"github.com/juancsr/platzi-grpc/repository"
	"github.com/juancsr/platzi-grpc/studentpb"
	"github.com/juancsr/platzi-grpc/testpb"
)

type TestServer struct {
	repo repository.Repository
	testpb.UnimplementedTestServiceServer
}

func NewTestServer(repo repository.Repository) *TestServer {
	return &TestServer{repo: repo}
}

func (s *TestServer) GetTest(ctx context.Context, req *testpb.GetTestRequest) (*testpb.Test, error) {
	test, err := s.repo.GetTest(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &testpb.Test{
		Id:   test.Id,
		Name: test.Name,
	}, nil
}

func (s *TestServer) SetTest(ctx context.Context, req *testpb.Test) (*testpb.SetTestResponse, error) {
	test := &models.Test{
		Id:   req.GetId(),
		Name: req.GetName(),
	}

	err := s.repo.SetTest(ctx, test)
	if err != nil {
		return nil, err
	}

	return &testpb.SetTestResponse{Id: test.Id, Name: test.Name}, nil
}

func (s *TestServer) SetQuestions(stream testpb.TestService_SetQuestionsServer) error {
	for {
		// se bloquea hasta que se reciba un mensaje del cliente
		msg, err := stream.Recv()
		if err == io.EOF {
			// ya no se envia más, responder la respuesta
			return stream.SendAndClose(&testpb.SetQuestionResponse{Ok: true})
		}

		if err != nil {
			return err
		}

		question := &models.Question{
			Id:       msg.GetId(),
			Answer:   msg.GetAnswer(),
			Question: msg.GetQuestion(),
			TestId:   msg.GetTestId(),
		}

		err = s.repo.SetQuestion(context.Background(), question)
		if err != nil {
			return stream.SendAndClose(&testpb.SetQuestionResponse{Ok: false})
		}
	}
}

func (s *TestServer) EnrollStudents(stream testpb.TestService_EnrollStudentsServer) error {
	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&testpb.SetQuestionResponse{Ok: true})
		}

		if err != nil {
			return err
		}

		enrollment := &models.Enrollment{
			StudentId: msg.GetStudentId(),
			TestId:    msg.GetTestId(),
		}

		err = s.repo.SetEnrollment(context.Background(), enrollment)

		if err != nil {
			return stream.SendAndClose(&testpb.SetQuestionResponse{Ok: false})
		}
	}
}

func (s *TestServer) GetStudentsPerTest(req *testpb.GetStudentsPerTestRequest,
	stream testpb.TestService_GetStudentsPerTestServer) error {
	stutends, err := s.repo.GetStudentsPerTest(context.Background(), req.GetTestId())
	if err != nil {
		return err
	}

	fmt.Println(stutends, req.GetTestId())

	for _, student := range stutends {
		studentStm := &studentpb.Student{
			Id:   student.Id,
			Name: student.Name,
			Age:  student.Age,
		}

		err := stream.Send(studentStm)
		// solo para probar
		time.Sleep(2 * time.Second)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *TestServer) TakeTest(stream testpb.TestService_TakeTestServer) error {
	questions, err := s.repo.GetQuestionsPerTest(context.Background(), "t1")
	if err != nil {
		return err
	}

	var (
		i               = 0
		currentQuestion = &models.Question{}
	)

	for {
		if i < len(questions) {
			currentQuestion = questions[i]
		}

		if i <= len(questions) {
			questionToSend := &testpb.Question{
				Id:       currentQuestion.Id,
				Question: currentQuestion.Question,
			}
			err := stream.Send(questionToSend)
			if err != nil {
				return err
			}
			i++
		}
		answer, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		log.Println(answer)
	}
}
