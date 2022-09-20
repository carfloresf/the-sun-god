// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package service

import (
	"github.com/carfloresf/the-sun-god/internal/model"
	"github.com/google/uuid"
	"sync"
)

// Ensure, that StorageMock does implement Storage.
// If this is not the case, regenerate this file with moq.
var _ Storage = &StorageMock{}

// StorageMock is a mock implementation of Storage.
//
//	func TestSomethingThatUsesStorage(t *testing.T) {
//
//		// make and configure a mocked Storage
//		mockedStorage := &StorageMock{
//			CreatePostFunc: func(post model.Post) (uuid.UUID, error) {
//				panic("mock out the CreatePost method")
//			},
//			GetPostsFunc: func(limit int, offset int) ([]model.Post, int, error) {
//				panic("mock out the GetPosts method")
//			},
//			GetPromotedPostsFunc: func(limit int) ([]model.Post, error) {
//				panic("mock out the GetPromotedPosts method")
//			},
//			UpdatePostScoreFunc: func(pID uuid.UUID, score int) error {
//				panic("mock out the UpdatePostScore method")
//			},
//		}
//
//		// use mockedStorage in code that requires Storage
//		// and then make assertions.
//
//	}
type StorageMock struct {
	// CreatePostFunc mocks the CreatePost method.
	CreatePostFunc func(post model.Post) (uuid.UUID, error)

	// GetPostsFunc mocks the GetPosts method.
	GetPostsFunc func(limit int, offset int) ([]model.Post, int, error)

	// GetPromotedPostsFunc mocks the GetPromotedPosts method.
	GetPromotedPostsFunc func(limit int) ([]model.Post, error)

	// UpdatePostScoreFunc mocks the UpdatePostScore method.
	UpdatePostScoreFunc func(pID uuid.UUID, score int) error

	// calls tracks calls to the methods.
	calls struct {
		// CreatePost holds details about calls to the CreatePost method.
		CreatePost []struct {
			// Post is the post argument value.
			Post model.Post
		}
		// GetPosts holds details about calls to the GetPosts method.
		GetPosts []struct {
			// Limit is the limit argument value.
			Limit int
			// Offset is the offset argument value.
			Offset int
		}
		// GetPromotedPosts holds details about calls to the GetPromotedPosts method.
		GetPromotedPosts []struct {
			// Limit is the limit argument value.
			Limit int
		}
		// UpdatePostScore holds details about calls to the UpdatePostScore method.
		UpdatePostScore []struct {
			// PID is the pID argument value.
			PID uuid.UUID
			// Score is the score argument value.
			Score int
		}
	}
	lockCreatePost       sync.RWMutex
	lockGetPosts         sync.RWMutex
	lockGetPromotedPosts sync.RWMutex
	lockUpdatePostScore  sync.RWMutex
}

// CreatePost calls CreatePostFunc.
func (mock *StorageMock) CreatePost(post model.Post) (uuid.UUID, error) {
	if mock.CreatePostFunc == nil {
		panic("StorageMock.CreatePostFunc: method is nil but Storage.CreatePost was just called")
	}
	callInfo := struct {
		Post model.Post
	}{
		Post: post,
	}
	mock.lockCreatePost.Lock()
	mock.calls.CreatePost = append(mock.calls.CreatePost, callInfo)
	mock.lockCreatePost.Unlock()
	return mock.CreatePostFunc(post)
}

// CreatePostCalls gets all the calls that were made to CreatePost.
// Check the length with:
//
//	len(mockedStorage.CreatePostCalls())
func (mock *StorageMock) CreatePostCalls() []struct {
	Post model.Post
} {
	var calls []struct {
		Post model.Post
	}
	mock.lockCreatePost.RLock()
	calls = mock.calls.CreatePost
	mock.lockCreatePost.RUnlock()
	return calls
}

// GetPosts calls GetPostsFunc.
func (mock *StorageMock) GetPosts(limit int, offset int) ([]model.Post, int, error) {
	if mock.GetPostsFunc == nil {
		panic("StorageMock.GetPostsFunc: method is nil but Storage.GetPosts was just called")
	}
	callInfo := struct {
		Limit  int
		Offset int
	}{
		Limit:  limit,
		Offset: offset,
	}
	mock.lockGetPosts.Lock()
	mock.calls.GetPosts = append(mock.calls.GetPosts, callInfo)
	mock.lockGetPosts.Unlock()
	return mock.GetPostsFunc(limit, offset)
}

// GetPostsCalls gets all the calls that were made to GetPosts.
// Check the length with:
//
//	len(mockedStorage.GetPostsCalls())
func (mock *StorageMock) GetPostsCalls() []struct {
	Limit  int
	Offset int
} {
	var calls []struct {
		Limit  int
		Offset int
	}
	mock.lockGetPosts.RLock()
	calls = mock.calls.GetPosts
	mock.lockGetPosts.RUnlock()
	return calls
}

// GetPromotedPosts calls GetPromotedPostsFunc.
func (mock *StorageMock) GetPromotedPosts(limit int) ([]model.Post, error) {
	if mock.GetPromotedPostsFunc == nil {
		panic("StorageMock.GetPromotedPostsFunc: method is nil but Storage.GetPromotedPosts was just called")
	}
	callInfo := struct {
		Limit int
	}{
		Limit: limit,
	}
	mock.lockGetPromotedPosts.Lock()
	mock.calls.GetPromotedPosts = append(mock.calls.GetPromotedPosts, callInfo)
	mock.lockGetPromotedPosts.Unlock()
	return mock.GetPromotedPostsFunc(limit)
}

// GetPromotedPostsCalls gets all the calls that were made to GetPromotedPosts.
// Check the length with:
//
//	len(mockedStorage.GetPromotedPostsCalls())
func (mock *StorageMock) GetPromotedPostsCalls() []struct {
	Limit int
} {
	var calls []struct {
		Limit int
	}
	mock.lockGetPromotedPosts.RLock()
	calls = mock.calls.GetPromotedPosts
	mock.lockGetPromotedPosts.RUnlock()
	return calls
}

// UpdatePostScore calls UpdatePostScoreFunc.
func (mock *StorageMock) UpdatePostScore(pID uuid.UUID, score int) error {
	if mock.UpdatePostScoreFunc == nil {
		panic("StorageMock.UpdatePostScoreFunc: method is nil but Storage.UpdatePostScore was just called")
	}
	callInfo := struct {
		PID   uuid.UUID
		Score int
	}{
		PID:   pID,
		Score: score,
	}
	mock.lockUpdatePostScore.Lock()
	mock.calls.UpdatePostScore = append(mock.calls.UpdatePostScore, callInfo)
	mock.lockUpdatePostScore.Unlock()
	return mock.UpdatePostScoreFunc(pID, score)
}

// UpdatePostScoreCalls gets all the calls that were made to UpdatePostScore.
// Check the length with:
//
//	len(mockedStorage.UpdatePostScoreCalls())
func (mock *StorageMock) UpdatePostScoreCalls() []struct {
	PID   uuid.UUID
	Score int
} {
	var calls []struct {
		PID   uuid.UUID
		Score int
	}
	mock.lockUpdatePostScore.RLock()
	calls = mock.calls.UpdatePostScore
	mock.lockUpdatePostScore.RUnlock()
	return calls
}
