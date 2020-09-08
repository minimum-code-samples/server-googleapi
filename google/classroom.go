package google

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"google.golang.org/api/classroom/v1"
	"google.golang.org/api/option"
)

// FetchAllCourses retrieves all Courses.
func FetchAllCourses(ctx context.Context, credentials []byte, token *oauth2.Token) ([]*classroom.Course, error) {
	srv, err := makeClassroomService(ctx, credentials, token)
	if err != nil {
		return nil, err
	}
	ipp := 3
	kourses := make([]*classroom.Course, 0)
	call := srv.Courses.List().PageSize(int64(ipp))
	call.Pages(ctx, func(resp *classroom.ListCoursesResponse) error {
		kourses = append(kourses, resp.Courses...)
		fmt.Printf("Page:\n%v\n", resp.Courses)
		return nil
	})
	return kourses, nil
}

// FetchCourses retrieves a single page of Courses.
func FetchCourses(ctx context.Context, credentials []byte, token *oauth2.Token) ([]*classroom.Course, error) {
	srv, err := makeClassroomService(ctx, credentials, token)
	if err != nil {
		return nil, err
	}
	call := srv.Courses.List()
	resp, err := call.Do()
	if err != nil {
		return nil, err
	}
	return resp.Courses, nil
}

// ReadClassroomNames extracts the names of the Courses from the supplied array.
func ReadClassroomNames(cr []*classroom.Course) []string {
	titles := make([]string, len(cr))
	for i, c := range cr {
		fmt.Printf("- ID: %s; Name: %s\n", c.Id, c.Name)
		titles[i] = c.Name
	}
	return titles
}

func makeClassroomService(ctx context.Context, credentials []byte, token *oauth2.Token) (*classroom.Service, error) {
	cfg, err := MakeConfig(credentials, ScopesWithClassroom())
	if err != nil {
		return nil, err
	}
	srv, err := classroom.NewService(ctx, option.WithTokenSource(cfg.TokenSource(ctx, token)))
	if err != nil {
		return nil, err
	}
	return srv, nil
}
