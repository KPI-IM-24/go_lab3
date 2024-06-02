package painter

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/shiny/screen"
)

func TestLoop_Post_Success(t *testing.T) {
	screenMock := new(mockScreen)
	textureMock := new(mockTexture)
	receiverMock := new(mockReceiver)

	texture := image.Pt(400, 400)
	screenMock.On("NewTexture", texture).Return(textureMock, nil)
	receiverMock.On("Update", textureMock).Return()

	loop := Loop{
		Receiver: receiverMock,
	}

	operation1 := new(mockOperation)
	textureMock.On("Bounds").Return(image.Rectangle{})
	operation1.On("Do", textureMock).Return(true)

	loop.Start(screenMock)
	loop.Post(operation1)
	time.Sleep(1 * time.Second)

	assert.Empty(t, loop.mq.ops)
	operation1.AssertCalled(t, "Do", textureMock)
	receiverMock.AssertCalled(t, "Update", textureMock)
	screenMock.AssertCalled(t, "NewTexture", image.Pt(400, 400))
}

func TestLoop_Post_Multiple_Success(t *testing.T) {
	screenMock := new(mockScreen)
	textureMock := new(mockTexture)
	receiverMock := new(mockReceiver)

	texture := image.Pt(400, 400)
	screenMock.On("NewTexture", texture).Return(textureMock, nil)
	receiverMock.On("Update", textureMock).Return()
	loop := Loop{
		Receiver: receiverMock,
	}

	loop.Start(screenMock)

	operation1 := new(mockOperation)
	operation2 := new(mockOperation)
	textureMock.On("Bounds").Return(image.Rectangle{})
	operation1.On("Do", textureMock).Return(true)
	operation2.On("Do", textureMock).Return(true)

	assert.Empty(t, loop.mq.ops)
	loop.Post(operation1)
	loop.Post(operation2)
	time.Sleep(1 * time.Second)
	assert.Empty(t, loop.mq.ops)

	operation1.AssertCalled(t, "Do", textureMock)
	operation2.AssertCalled(t, "Do", textureMock)
	receiverMock.AssertCalled(t, "Update", textureMock)
	screenMock.AssertCalled(t, "NewTexture", image.Pt(400, 400))
}

func TestLoop_Post_Failure(t *testing.T) {
	screenMock := new(mockScreen)
	textureMock := new(mockTexture)
	receiverMock := new(mockReceiver)

	texture := image.Pt(400, 400)
	screenMock.On("NewTexture", texture).Return(textureMock, nil)
	receiverMock.On("Update", textureMock).Return()
	loop := Loop{
		Receiver: receiverMock,
	}

	loop.Start(screenMock)
	operation1 := new(mockOperation)
	textureMock.On("Bounds").Return(image.Rectangle{})
	operation1.On("Do", textureMock).Return(false)

	assert.Empty(t, loop.mq.ops)
	loop.Post(operation1)
	time.Sleep(1 * time.Second)
	assert.Empty(t, loop.mq.ops)

	operation1.AssertCalled(t, "Do", textureMock)
	receiverMock.AssertNotCalled(t, "Update", textureMock)
	screenMock.AssertCalled(t, "NewTexture", image.Pt(400, 400))
}

type mockScreen struct {
	mock.Mock
}

func (m *mockScreen) NewBuffer(size image.Point) (screen.Buffer, error) {
	return nil, nil
}

func (m *mockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	args := m.Called(size)
	return args.Get(0).(screen.Texture), args.Error(1)
}

func (m *mockScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	return nil, nil
}

type mockTexture struct {
	mock.Mock
}

func (m *mockTexture) Release() {
	m.Called()
}

func (m *mockTexture) Size() image.Point {
	args := m.Called()
	return args.Get(0).(image.Point)
}

func (m *mockTexture) Bounds() image.Rectangle {
	args := m.Called()
	return args.Get(0).(image.Rectangle)
}

func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {
	m.Called(dp, src, sr)
}

func (m *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	m.Called(dr, src, op)
}

type mockReceiver struct {
	mock.Mock
}

func (m *mockReceiver) Update(texture screen.Texture) {
	m.Called(texture)
}

type mockOperation struct {
	mock.Mock
}

func (m *mockOperation) Do(t screen.Texture) bool {
	args := m.Called(t)
	return args.Bool(0)
}
