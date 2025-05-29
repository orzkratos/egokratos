package example2

import (
	"fmt"

	"github.com/orzkratos/errkratos"
)

type Class struct {
	Name string
}

type ClassStudentsScores struct {
	Class          *Class
	StudentsScores []*StudentScores
	AvgScore       float64
	Erk            *errkratos.Erk
}

type Student struct {
	Name string
}

type StudentScores struct {
	Student  *Student
	Scores   []*SubjectScore
	AvgScore float64
	Erk      *errkratos.Erk
}

type Subject struct {
	Name string
}

type SubjectScore struct {
	Subject *Subject
	Score   int
	Erk     *errkratos.Erk
}

func NewClasses(classCount int) []*Class {
	classes := make([]*Class, 0, classCount)
	for idx := 0; idx < classCount; idx++ {
		classes = append(classes, &Class{
			Name: fmt.Sprintf("class(%d)", idx),
		})
	}
	return classes
}

func NewStudents(studentCount int) []*Student {
	var students = make([]*Student, 0, studentCount)
	for idx := 0; idx < studentCount; idx++ {
		students = append(students, &Student{
			Name: fmt.Sprintf("student(%d)", idx),
		})
	}
	return students
}

func NewSubjects(subjectCount int) []*Subject {
	var subjects = make([]*Subject, 0, subjectCount)
	for idx := 0; idx < subjectCount; idx++ {
		subjects = append(subjects, &Subject{
			Name: fmt.Sprintf("subject(%d)", idx),
		})
	}
	return subjects
}
