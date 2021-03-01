package integration

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/jmoiron/sqlx"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const dateLayout = "2006-01-02 15:04"

type Suite struct {
	suite.Suite

	db          *sqlx.DB
	grpcConn    *grpc.ClientConn
	eventClient pb.EventServiceClient
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupSuite() {
	var err error
	s.grpcConn, err = grpc.Dial("calendar_api:8082", grpc.WithInsecure())
	s.Require().NoError(err)

	s.eventClient = pb.NewEventServiceClient(s.grpcConn)

	addr := "calendar_user:calendar_pass@tcp(calendar_db:3306)/calendar?parseTime=true"
	s.db, err = sqlx.Connect("mysql", addr)
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.DangerousSkipTestDatabaseCheck(),
		testfixtures.Database(s.db.DB),
		testfixtures.Dialect("mysql"),
		testfixtures.Directory("fixtures"),
	)
	s.Require().NoError(err)

	err = fixtures.Load()
	s.Require().NoError(err)
}

func (s *Suite) TearDownSuite() {
	s.Require().NoError(s.db.Close())
	s.Require().NoError(s.grpcConn.Close())
}

func (s *Suite) TestGetEventByID() {
	s.Run("ok", func() {
		resp, err := s.eventClient.GetEventByID(context.Background(), &pb.GetEventByIDRequest{
			Id: 1,
		})
		s.Require().NoError(err)

		s.Equal(int64(1), resp.Event.Id)
		s.Equal("title here", resp.Event.Title)
		s.Equal("description here", resp.Event.Description)
		s.Equal(int64(1), resp.Event.UserId)
		s.Equal("2099-04-01 10:00", resp.Event.StartDate.AsTime().Format(dateLayout))
		s.Equal("2099-04-01 10:10", resp.Event.EndDate.AsTime().Format(dateLayout))
		s.Equal(int32(0), resp.Event.IsNotified)
	})

	s.Run("not existing event", func() {
		resp, err := s.eventClient.GetEventByID(context.Background(), &pb.GetEventByIDRequest{
			Id: 100500,
		})
		s.Require().Equal(codes.NotFound, status.Code(err))
		s.Require().Nil(resp)
	})
}

func (s *Suite) TestCreateEvent() {
	s.Run("ok", func() {
		pbEvent := &pb.Event{
			Title:            "title 999",
			Description:      "descr 999",
			UserId:           111,
			StartDate:        timestamppb.New(time.Now().AddDate(1, 0, 0)),
			EndDate:          timestamppb.New(time.Now().AddDate(1, 0, 1)),
			NotificationDate: timestamppb.New(time.Now().AddDate(1, 0, 2)),
		}

		resp, err := s.eventClient.CreateEvent(context.Background(), &pb.CreateEventRequest{
			Event: pbEvent,
		})
		s.Require().NoError(err)

		e, err := s.fetchEvent(resp.InsertedId)
		s.Require().NoError(err)
		s.Require().Equal(storage.EventID(resp.InsertedId), e.ID)
		s.compareEvents(pbEvent, e)
	})

	s.Run("busy date error", func() {
		sdate, err := time.Parse(dateLayout, "2099-04-01 10:00")
		s.Require().NoError(err)
		pbEvent := &pb.Event{
			Title:       "new title",
			Description: "new descr",
			UserId:      1,
			StartDate:   timestamppb.New(sdate),
		}

		resp, err := s.eventClient.CreateEvent(context.Background(), &pb.CreateEventRequest{
			Event: pbEvent,
		})
		s.Require().Error(err)
		s.Require().Nil(resp)
	})
}

func (s *Suite) TestDeleteEventByID() {
	s.Run("ok", func() {
		var id int64 = 2
		resp, err := s.eventClient.DeleteEvent(context.Background(), &pb.DeleteEventRequest{
			Id: id,
		})
		s.Require().NoError(err)
		s.Require().Equal(int64(1), resp.Affected)

		_, err = s.fetchEvent(id)
		s.Require().True(errors.Is(err, sql.ErrNoRows))
	})

	s.Run("not existing event", func() {
		resp, err := s.eventClient.DeleteEvent(context.Background(), &pb.DeleteEventRequest{
			Id: 100500,
		})
		s.Require().NoError(err)
		s.Require().Equal(int64(0), resp.Affected)
	})
}

func (s *Suite) TestUpdateEvent() {
	s.Run("ok", func() {
		var id int64 = 3
		pbEvent := &pb.Event{
			Id:               id,
			Title:            "updated title",
			Description:      "updated descr",
			UserId:           101,
			StartDate:        timestamppb.New(time.Now().AddDate(1, 0, 0)),
			EndDate:          timestamppb.New(time.Now().AddDate(1, 0, 1)),
			NotificationDate: timestamppb.New(time.Now().AddDate(1, 0, 2)),
		}

		resp, err := s.eventClient.UpdateEvent(context.Background(), &pb.UpdateEventRequest{
			Id:    id,
			Event: pbEvent,
		})
		s.Require().NoError(err)
		s.Require().Equal(int64(1), resp.Affected)

		e, err := s.fetchEvent(id)
		s.Require().NoError(err)
		s.compareEvents(pbEvent, e)
	})

	s.Run("not existing event", func() {
		var id int64 = 303500
		resp, err := s.eventClient.UpdateEvent(context.Background(), &pb.UpdateEventRequest{
			Id: id,
			Event: &pb.Event{
				Id:          id,
				Title:       "updated title",
				Description: "updated descr",
			},
		})
		s.Require().NoError(err)
		s.Require().Equal(int64(0), resp.Affected)
	})

	s.Run("busy date error", func() {
		var id int64 = 3
		sdate, err := time.Parse(dateLayout, "2099-04-01 10:00")
		s.Require().NoError(err)

		resp, err := s.eventClient.UpdateEvent(context.Background(), &pb.UpdateEventRequest{
			Id: id,
			Event: &pb.Event{
				Id:        id,
				UserId:    1,
				StartDate: timestamppb.New(sdate),
			},
		})
		s.Require().Error(err)
		s.Require().Nil(resp)
	})
}

func (s *Suite) TestGetUserDayEvents() {
	s.Run("ok", func() {
		ids := map[int64]struct{}{
			7: {},
			8: {},
		}

		date, err := time.Parse(dateLayout, "2100-05-05 00:00")
		s.Require().NoError(err)

		resp, err := s.eventClient.GetUserDayEvents(context.Background(), &pb.UserPeriodEventRequest{
			UserID: 500,
			Date:   timestamppb.New(date),
		})
		s.Require().NoError(err)
		s.Require().Equal(len(ids), len(resp.Events))
		for _, e := range resp.Events {
			s.Require().Contains(ids, e.Id)
		}
	})

	s.Run("no events", func() {
		date, err := time.Parse(dateLayout, "2000-05-05 00:00")
		s.Require().NoError(err)

		resp, err := s.eventClient.GetUserDayEvents(context.Background(), &pb.UserPeriodEventRequest{
			UserID: 500,
			Date:   timestamppb.New(date),
		})
		s.Require().NoError(err)
		s.Require().Equal(0, len(resp.Events))
	})
}

func (s *Suite) TestGetUserWeekEvents() {
	s.Run("ok", func() {
		ids := map[int64]struct{}{
			7: {},
			8: {},
			9: {},
		}

		date, err := time.Parse(dateLayout, "2100-05-05 00:00")
		s.Require().NoError(err)

		resp, err := s.eventClient.GetUserWeekEvents(context.Background(), &pb.UserPeriodEventRequest{
			UserID: 500,
			Date:   timestamppb.New(date),
		})
		s.Require().NoError(err)
		s.Require().Equal(len(ids), len(resp.Events))
		for _, e := range resp.Events {
			s.Require().Contains(ids, e.Id)
		}
	})

	s.Run("no events", func() {
		date, err := time.Parse(dateLayout, "2000-05-05 00:00")
		s.Require().NoError(err)

		resp, err := s.eventClient.GetUserWeekEvents(context.Background(), &pb.UserPeriodEventRequest{
			UserID: 500,
			Date:   timestamppb.New(date),
		})
		s.Require().NoError(err)
		s.Require().Equal(0, len(resp.Events))
	})
}

func (s *Suite) TestGetUserMonthEvents() {
	s.Run("ok", func() {
		ids := map[int64]struct{}{
			7:  {},
			8:  {},
			9:  {},
			10: {},
		}

		date, err := time.Parse(dateLayout, "2100-05-05 00:00")
		s.Require().NoError(err)

		resp, err := s.eventClient.GetUserMonthEvents(context.Background(), &pb.UserPeriodEventRequest{
			UserID: 500,
			Date:   timestamppb.New(date),
		})
		s.Require().NoError(err)
		s.Require().Equal(len(ids), len(resp.Events))
		for _, e := range resp.Events {
			s.Require().Contains(ids, e.Id)
		}
	})

	s.Run("no events", func() {
		date, err := time.Parse(dateLayout, "2000-05-05 00:00")
		s.Require().NoError(err)

		resp, err := s.eventClient.GetUserMonthEvents(context.Background(), &pb.UserPeriodEventRequest{
			UserID: 500,
			Date:   timestamppb.New(date),
		})
		s.Require().NoError(err)
		s.Require().Equal(0, len(resp.Events))
	})
}

func (s *Suite) TestSender() {
	time.Sleep(10 * time.Second)

	s.Run("notification", func() {
		e, err := s.fetchEvent(6)
		s.Require().NoError(err)
		s.Require().Equal(byte(1), e.IsNotified)
	})

	s.Run("delete old event", func() {
		_, err := s.fetchEvent(4)
		s.Require().True(errors.Is(err, sql.ErrNoRows))
	})
}

func (s *Suite) compareEvents(pbEvent *pb.Event, e *storage.Event) {
	s.Require().Equal(pbEvent.Description, e.Description)
	s.Require().Equal(storage.UserID(pbEvent.UserId), e.UserID)
	s.Require().Equal(pbEvent.StartDate.AsTime().Format(dateLayout), e.StartDate.Format(dateLayout))
	s.Require().Equal(pbEvent.EndDate.AsTime().Format(dateLayout), e.EndDate.Format(dateLayout))
	s.Require().Equal(pbEvent.NotificationDate.AsTime().Format(dateLayout), e.NotificationDate.Format(dateLayout))
	s.Require().Equal(pbEvent.IsNotified, int32(e.IsNotified))
}

func (s *Suite) fetchEvent(id int64) (*storage.Event, error) {
	event := new(storage.Event)
	err := s.db.
		QueryRowx("SELECT * FROM event WHERE id = ?", id).
		StructScan(event)

	return event, err
}
