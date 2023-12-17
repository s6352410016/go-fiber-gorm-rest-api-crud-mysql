package handlers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/s6352410016/go-fiber-gorm-rest-api-crud-mysql/database"
	"github.com/s6352410016/go-fiber-gorm-rest-api-crud-mysql/models"
)

func Create(c *fiber.Ctx) error {
	employee := new(models.Employee)
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "Error While File Upload",
		})
	}

	fileExt := filepath.Ext(file.Filename)
	// filter file extension
	allowExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}
	if !allowExts[fileExt] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "Invalid File Extension",
		})
	}

	// change file name
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExt)

	empName := c.FormValue("empName")
	empAddress := c.FormValue("empAddress")
	empTel := c.FormValue("empTel")
	empSalary := c.FormValue("empSalary")
	empPhoto := newFileName

	if empName == "" || empAddress == "" || empTel == "" || empSalary == "" || empPhoto == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "Invalid Request Data",
		})
	}

	empSalaryFloat, err := strconv.ParseFloat(empSalary, 32)

	if err != nil {
		fmt.Println("Error While Parse String To Float\n", err.Error())
	}

	employee.EmpName = empName
	employee.EmpAddress = empAddress
	employee.EmpTel = empTel
	employee.EmpSalary = float32(empSalaryFloat)
	employee.EmpPhoto = empPhoto

	database.DB.Create(&employee)

	// save file to directory
	if err := c.SaveFile(file, "./public/"+newFileName); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "Server Error When Save File",
		})
	}

	return c.Status(fiber.StatusOK).JSON(employee)
}

func GetAll(c *fiber.Ctx) error {
	var employees []models.Employee
	database.DB.Find(&employees)
	return c.Status(fiber.StatusOK).JSON(employees)
}

func GetById(c *fiber.Ctx) error {
	employeeId := c.Params("id")
	var employee models.Employee
	database.DB.First(&employee, employeeId)

	if employee.EmpName == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "Employee Not Found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(employee)
}

func Update(c *fiber.Ctx) error {
	employeeId := c.Params("id")
	var employee models.Employee
	database.DB.First(&employee, employeeId)

	if employee.EmpName == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "Employee Not Found",
		})
	}

	file, _ := c.FormFile("image")

	// upload file
	if file != nil {
		fileExt := filepath.Ext(file.Filename)
		allowExts := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".webp": true,
		}

		if !allowExts[fileExt] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"msg": "Invalid File Extension",
			})
		}

		newFileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExt)

		if c.FormValue("empSalary") != "" {
			empSalaryFloat, err := strconv.ParseFloat(c.FormValue("empSalary"), 32)
			if err != nil {
				fmt.Println("Error While Parse String To Float\n", err.Error())
			}

			empName := c.FormValue("empName")
			empAddress := c.FormValue("empAddress")
			empTel := c.FormValue("empTel")
			empSalary := float32(empSalaryFloat)
			empPhoto := newFileName

			if empName == "" || empAddress == "" || empTel == "" || empSalary == 0 || empPhoto == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"msg": "Invalid Request Data",
				})
			}

			// delete current employee photo
			os.Remove("./public/" + employee.EmpPhoto)

			employee.EmpName = empName
			employee.EmpAddress = empAddress
			employee.EmpTel = empTel
			employee.EmpSalary = empSalary
			employee.EmpPhoto = empPhoto

			database.DB.Save(&employee)

			if err := c.SaveFile(file, "./public/"+empPhoto); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"msg": "Server Error When Save File",
				})
			}

			return c.Status(fiber.StatusOK).JSON(employee)
		}
	}

	// not upload file
	if c.FormValue("empSalary") != "" {
		empSalaryFloat, err := strconv.ParseFloat(c.FormValue("empSalary"), 32)
		if err != nil {
			fmt.Println("Error While Parse String To Float\n", err.Error())
		}

		empName := c.FormValue("empName")
		empAddress := c.FormValue("empAddress")
		empTel := c.FormValue("empTel")
		empSalary := float32(empSalaryFloat)

		if empName == "" || empAddress == "" || empTel == "" || empSalary == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"msg": "Invalid Request Data",
			})
		}

		employee.EmpName = empName
		employee.EmpAddress = empAddress
		employee.EmpTel = empTel
		employee.EmpSalary = empSalary

		database.DB.Save(&employee)
	}

	return c.Status(fiber.StatusOK).JSON(employee)
}

func Delete(c *fiber.Ctx) error {
	employeeId := c.Params("id")
	var employee models.Employee
	database.DB.First(&employee, employeeId)

	if employee.EmpName == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "Employee Not Found",
		})
	}

	err := os.Remove("./public/" + employee.EmpPhoto)
	if err != nil {
		fmt.Println("Error While Delete Photo\n", err.Error())
	}

	database.DB.Delete(&employee)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": "Employee Deleted Successfully",
	})
}

func GetImage(c *fiber.Ctx) error {
	fileName := c.Params("filename")

	// read file image from directory
	filePath := "./public/" + fileName
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error While Open Directory\n", err.Error())
	}
	defer file.Close()

	// read file data from file image
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "Image Not Found",
		})
	}

	// send image to client
	return c.Send(fileData)
}
