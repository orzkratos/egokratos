package example2_test

import (
	"context"
	"math/rand/v2"
	"testing"

	"github.com/orzkratos/egokratos"
	"github.com/orzkratos/egokratos/erkgroup"
	"github.com/orzkratos/egokratos/internal/errorspb"
	"github.com/orzkratos/egokratos/internal/examples/example2"
	"github.com/orzkratos/errkratos"
	"github.com/orzkratos/errkratos/must/erkmust"
	"github.com/yyle88/neatjson/neatjsons"
)

// TestRun demonstrates three-stage nested batch processing: classes -> students -> subjects
// Shows multi-stage task processing with score aggregation at each stage
//
// TestRun 演示三阶段嵌套批量处理：班级 -> 学生 -> 科目
// 展示多阶段任务处理，每个阶段都有分数聚合
func TestRun(t *testing.T) {
	ctx := context.Background()
	classes := example2.NewClasses(5)
	taskResults := processClasses(ctx, classes)
	classesStudentsScores := taskResults.Flatten(func(arg *example2.Class, erk *errkratos.Erk) *example2.ClassStudentsScores {
		return &example2.ClassStudentsScores{
			Class:          arg,
			StudentsScores: nil,
			Erk:            erk,
		}
	})
	t.Log(neatjsons.S(classesStudentsScores))
}

func processClasses(ctx context.Context, classes []*example2.Class) egokratos.Tasks[*example2.Class, *example2.ClassStudentsScores] {
	taskBatch := egokratos.NewTaskBatch[*example2.Class, *example2.ClassStudentsScores](classes)
	taskBatch.SetGlide(true)
	taskBatch.SetWaCtx(func(err error) *errkratos.Erk {
		return errorspb.ErrorWrongContext("wrong-ctx-can-not-invoke-process-class-func. error=%v", err)
	})
	ego := erkgroup.NewGroup(ctx)
	ego.SetLimit(3)
	taskBatch.EgoRun(ego, processClassFunc)
	erkmust.Done(ego.Wait())
	return taskBatch.Tasks
}

func processClassFunc(ctx context.Context, arg *example2.Class) (*example2.ClassStudentsScores, *errkratos.Erk) {
	if rand.IntN(100) < 30 {
		return nil, errorspb.ErrorServerDbError("wrong-db")
	}
	studentCount := 1 + rand.IntN(5)
	students := example2.NewStudents(studentCount)
	taskResults := processStudents(ctx, students)
	studentsScores := taskResults.Flatten(func(arg *example2.Student, erk *errkratos.Erk) *example2.StudentScores {
		return &example2.StudentScores{
			Student:  arg,
			Scores:   nil,
			AvgScore: 0,
			Erk:      erk,
		}
	})

	okCnt := 0
	okSum := float64(0)
	for _, studentScores := range studentsScores {
		if studentScores.Erk != nil {
			continue
		}
		okCnt++
		okSum += studentScores.AvgScore
	}
	avgScore := float64(0)
	if okCnt > 0 {
		avgScore = okSum / float64(okCnt)
	}

	return &example2.ClassStudentsScores{
		Class:          arg,
		StudentsScores: studentsScores,
		AvgScore:       avgScore,
		Erk:            nil,
	}, nil
}

func processStudents(ctx context.Context, students []*example2.Student) egokratos.Tasks[*example2.Student, *example2.StudentScores] {
	taskBatch := egokratos.NewTaskBatch[*example2.Student, *example2.StudentScores](students)
	taskBatch.SetGlide(true)
	taskBatch.SetWaCtx(func(err error) *errkratos.Erk {
		return errorspb.ErrorWrongContext("wrong-ctx-can-not-invoke-process-student-func. error=%v", err)
	})
	ego := erkgroup.NewGroup(ctx)
	ego.SetLimit(2)
	taskBatch.EgoRun(ego, processStudentFunc)
	erkmust.Done(ego.Wait())
	return taskBatch.Tasks
}

func processStudentFunc(ctx context.Context, arg *example2.Student) (*example2.StudentScores, *errkratos.Erk) {
	if rand.IntN(100) < 30 {
		return nil, errorspb.ErrorServerDbError("wrong-db")
	}
	subjectCount := 1 + rand.IntN(2)
	subjects := example2.NewSubjects(subjectCount)

	taskResults := processSubjects(ctx, subjects)
	scores := taskResults.Flatten(func(arg *example2.Subject, erk *errkratos.Erk) *example2.SubjectScore {
		return &example2.SubjectScore{
			Subject: arg,
			Score:   0,
			Erk:     erk,
		}
	})

	okCnt := 0
	okSum := float64(0)
	for _, score := range scores {
		if score.Erk != nil {
			continue
		}
		okCnt++
		okSum += float64(score.Score)
	}
	avgScore := float64(0)
	if okCnt > 0 {
		avgScore = okSum / float64(okCnt)
	}

	return &example2.StudentScores{
		Student:  arg,
		Scores:   scores,
		AvgScore: avgScore,
		Erk:      nil,
	}, nil
}

func processSubjects(ctx context.Context, subjects []*example2.Subject) egokratos.Tasks[*example2.Subject, *example2.SubjectScore] {
	taskBatch := egokratos.NewTaskBatch[*example2.Subject, *example2.SubjectScore](subjects)
	taskBatch.SetGlide(true)
	taskBatch.SetWaCtx(func(err error) *errkratos.Erk {
		return errorspb.ErrorWrongContext("wrong-ctx-can-not-invoke-process-subject-func. error=%v", err)
	})
	ego := erkgroup.NewGroup(ctx)
	ego.SetLimit(2)
	taskBatch.EgoRun(ego, processSubjectFunc)
	erkmust.Done(ego.Wait())
	return taskBatch.Tasks
}

func processSubjectFunc(ctx context.Context, arg *example2.Subject) (*example2.SubjectScore, *errkratos.Erk) {
	if rand.IntN(100) < 30 {
		return nil, errorspb.ErrorServerDbError("wrong-db")
	}
	return &example2.SubjectScore{
		Subject: arg,
		Score:   rand.IntN(100),
		Erk:     nil,
	}, nil
}
