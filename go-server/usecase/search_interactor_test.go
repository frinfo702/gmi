package usecase_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/frinfo702/rustysearch/domain/model"
	"github.com/frinfo702/rustysearch/usecase"
	"github.com/frinfo702/rustysearch/usecase/mock"
	"go.uber.org/mock/gomock"
)

func TestSearchInteractor_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// モックの Searcher を生成
	mockSearcher := mock.NewMockSearcher(ctrl)

	query := model.SearchQuery{Path: ".", Query: "test", Fuzzy: false}
	expectedResults := []model.SearchResult{
		{FilePath: "file1.txt", LineNumber: 1, LineText: "this is a test"},
	}

	// Search メソッドが呼ばれることを期待
	mockSearcher.EXPECT().Search(query).Return(expectedResults, nil)

	interactor := usecase.NewSearchInteractor(mockSearcher)
	results, err := interactor.Execute(query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(results, expectedResults) {
		t.Fatalf("expected %v, got %v", expectedResults, results)
	}
}

func TestSearchInteractor_Execute_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSearcher := mock.NewMockSearcher(ctrl)

	query := model.SearchQuery{Path: ".", Query: "fail", Fuzzy: false}
	expectedErr := errors.New("search failed")
	mockSearcher.EXPECT().Search(query).Return(nil, expectedErr)

	interactor := usecase.NewSearchInteractor(mockSearcher)
	_, err := interactor.Execute(query)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}
