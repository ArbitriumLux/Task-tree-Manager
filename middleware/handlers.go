package middleware

import (
	"TasksManager/models"
	"TasksManager/server"
	"html/template"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

//IndexPage assembly Index Page
func IndexPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/index.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

//CreateTask
func CreateTask(w http.ResponseWriter, r *http.Request) {
	server.DataBaseConnection()
	defer server.Db.Close()
	var task models.Tasks
	task.TaskName = r.FormValue("TaskName")
	if err := server.Db.Create(&task).Error; err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	http.Redirect(w, r, "/", 301)
	log.Info(task.TaskName, " created")
}

//CreateMiniTask
func CreateMiniTask(w http.ResponseWriter, r *http.Request) {
	server.DataBaseConnection()
	defer server.Db.Close()
	var task models.Tasks
	if err := server.Db.Where("Task_Name = ?", r.FormValue("TaskName")).Find(&task).Error; err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	var miniTask models.MiniTasks
	miniTask.MiniTaskName = r.FormValue("MiniTaskName")
	task.MiniTasks = append(task.MiniTasks, miniTask)
	if err := server.Db.Save(&task).Error; err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	http.Redirect(w, r, "/", 301)
	log.Info(miniTask.MiniTaskName, " for ", task.TaskName, " created")
}

//CreateLaborCost
func CreateLaborCost(w http.ResponseWriter, r *http.Request) {
	server.DataBaseConnection()
	defer server.Db.Close()
	var miniTask models.MiniTasks
	if err := server.Db.Where("Mini_Task_Name = ?", r.FormValue("MiniTaskName")).Find(&miniTask).Error; err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	var laborCost models.LaborCosts
	var err error
	laborCost.Cost, err = strconv.Atoi(r.FormValue("Cost"))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	miniTask.LaborCosts = append(miniTask.LaborCosts, laborCost)
	if err := server.Db.Save(&miniTask).Error; err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	http.Redirect(w, r, "/", 301)
	log.Info(laborCost.Cost, " for ", miniTask.MiniTaskName, " created")
}

type GetTreeResult struct {
	models.Tasks
	Sum     int
	SumSR   int
	SumMini []int
	SumsSR  []int
}

//GetTree fills the template with data by Task name
func GetTree(w http.ResponseWriter, r *http.Request) {
	server.DataBaseConnection()
	defer server.Db.Close()
	tmpl, err := template.ParseFiles("static/tree.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	var task models.Tasks
	params := mux.Vars(r)
	if err := server.Db.Preload("MiniTasks.LaborCosts").Find(&task, "Task_Name = ?", params["TaskName"]).Error; err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	sum, sumSR, sumMini, sumsSR := Sum(task)
	result := GetTreeResult{
		Tasks:   task,
		Sum:     sum,
		SumSR:   sumSR,
		SumMini: sumMini,
		SumsSR:  sumsSR,
	}
	if err := tmpl.Execute(w, result); err != nil {
		log.Print(err.Error())
		return
	}
}

//Sum calculates the total time spent and the average time for all subordinate elements
func Sum(task models.Tasks) (int, int, []int, []int) {
	var sum int
	var sumSR int
	sums := make([]int, len(task.MiniTasks))
	sumsSR := make([]int, len(task.MiniTasks))
	var length int
	var mutex sync.Mutex
	var wg sync.WaitGroup
	for index, miniTasks := range task.MiniTasks {
		wg.Add(1)
		length += len(miniTasks.LaborCosts)
		go func(index int, miniTasks models.MiniTasks) {
			defer wg.Done()
			sumMini := 0
			for _, LaborCosts := range miniTasks.LaborCosts {
				mutex.Lock()
				sum += LaborCosts.Cost
				mutex.Unlock()
				sumMini += LaborCosts.Cost
			}
			if length == 0 {
				return
			}
			if len(miniTasks.LaborCosts) == 0 {
				return
			}
			sumSR = sum / length
			sums[index] = sumMini
			sumsSR[index] = sumMini / len(miniTasks.LaborCosts)
		}(index, miniTasks)
	}
	wg.Wait()
	return sum, sumSR, sums, sumsSR
}

//UpdateMiniTask reassigns a Minitask to another parent
func UpdateMiniTask(w http.ResponseWriter, r *http.Request) {
	server.DataBaseConnection()
	defer server.Db.Close()
	http.Redirect(w, r, r.Header.Get("Referer"), 301)
	var task models.Tasks
	if err := server.Db.Select("id").Where("Task_Name =?", r.FormValue("TaskName")).Find(&task).Error; err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	var miniTask models.MiniTasks
	if err := server.Db.Where("Mini_Task_Name = ?", r.FormValue("MiniTaskName")).Find(&miniTask).Error; err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	miniTask.TasksID = task.ID
	if err := server.Db.Save(&miniTask).Error; err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	log.Info(miniTask.MiniTaskName, " reassigned to ", r.FormValue("TaskName"))
}

//UpdateLaborCost reassigns a Labor Cost to another parent
func UpdateLaborCost(w http.ResponseWriter, r *http.Request) {
	server.DataBaseConnection()
	defer server.Db.Close()
	http.Redirect(w, r, r.Header.Get("Referer"), 301)
	var miniTask models.MiniTasks
	if err := server.Db.Select("id").Where("Mini_Task_Name =?", r.FormValue("MiniTaskName")).Find(&miniTask).Error; err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	var laborCost models.LaborCosts
	if err := server.Db.Where("ID = ?", r.FormValue("ID")).Find(&laborCost).Error; err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	laborCost.MiniTasksID = miniTask.ID
	if err := server.Db.Save(&laborCost).Error; err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	log.Info(laborCost.Cost, " reassigned to ", r.FormValue("MiniTaskName"))
}

//DeleteTask
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, r.Header.Get("Referer"), 301)
	server.DataBaseConnection()
	defer server.Db.Close()
	emp, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err := server.Db.Delete(&models.Tasks{}, emp).Error; err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	log.Info("TaskID= ", emp, " deleted")
}

//DeleteMiniTask
func DeleteMiniTask(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, r.Header.Get("Referer"), 301)
	server.DataBaseConnection()
	defer server.Db.Close()
	emp, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err := server.Db.Delete(&models.MiniTasks{}, emp).Error; err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	log.Info("MinitaskID= ", emp, " deleted")
}

//DeleteLaborCost
func DeleteLaborCost(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, r.Header.Get("Referer"), 301)
	server.DataBaseConnection()
	defer server.Db.Close()
	emp, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err := server.Db.Delete(&models.LaborCosts{}, emp).Error; err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	log.Info("Labor_CostID= ", emp, " deleted")
}
