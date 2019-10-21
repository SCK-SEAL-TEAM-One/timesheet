package repository

import (
	"timesheet/internal/model"

	"github.com/jmoiron/sqlx"
)

type TimesheetRepositoryGateways interface {
	GetSummary(year, month int) ([]model.TransactionTimesheet, error)
	GetMemberByID(memberID string) ([]model.Member, error)
	GetIncomes(memberID string, year, month int) ([]model.Incomes, error)
	CreateIncome(year, month int, memberID string, income model.Incomes) error
}

type TimesheetRepository struct {
	DatabaseConnection *sqlx.DB
}

func (repository TimesheetRepository) GetSummary(year, month int) ([]model.TransactionTimesheet, error) {
	var transactionTimesheetList []model.TransactionTimesheet
	query := `SELECT * FROM transactions WHERE year = ? AND month = ? ORDER BY member_id ASC, company DESC`
	err := repository.DatabaseConnection.Select(&transactionTimesheetList, query, year, month)
	if err != nil {
		return []model.TransactionTimesheet{}, err
	}
	return transactionTimesheetList, nil
}

func (repository TimesheetRepository) CreateIncome(year, month int, memberID string, income model.Incomes) error {
	query := `INSERT INTO incomes (member_id, month, year, day, start_time_am_hours, 
		start_time_am_minutes, start_time_am_seconds, end_time_am_hours, end_time_am_minutes, 
		end_time_am_seconds, start_time_pm_hours, start_time_pm_minutes, start_time_pm_seconds, 
		end_time_pm_hours, end_time_pm_minutes, end_time_pm_seconds, overtime, total_hours_hours, 
		total_hours_minutes, total_hours_seconds, coaching_customer_charging, coaching_payment_rate, 
		training_wage, other_wage, company, description) 
		VALUES ( ? , ? , ? , ? , ? , ? , ? , ? , ? , ? , ? , ? , ?, ? , ? , 
			?, ? , ? , ?, ? , ? , ?, ? , ? , ?, ? )`
	transaction := repository.DatabaseConnection.MustBegin()
	transaction.MustExec(query, memberID, month, year, income.Day, income.StartTimeAMHours, income.StartTimeAMMinutes, income.StartTimeAMSeconds, income.EndTimeAMHours, income.EndTimeAMMinutes, income.EndTimeAMSeconds, income.StartTimePMHours, income.StartTimePMMinutes, income.StartTimePMSeconds, income.EndTimePMHours, income.EndTimePMMinutes, income.EndTimePMSeconds, income.Overtime, income.TotalHoursHours, income.TotalHoursMinutes, income.TotalHoursSeconds, income.CoachingCustomerCharging, income.CoachingPaymentRate, income.TrainingWage, income.OtherWage, income.Company, income.Description)
	err := transaction.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (repository TimesheetRepository) GetMemberByID(memberID string) ([]model.Member, error) {
	var memberList []model.Member
	query := `SELECT * FROM timesheet.members WHERE member_id = ?`
	err := repository.DatabaseConnection.Select(&memberList, query, memberID)
	if err != nil {
		return []model.Member{}, err
	}
	return memberList, nil
}

func (repository TimesheetRepository) GetIncomes(memberID string, year, month int) ([]model.Incomes, error) {
	var incomeList []model.Incomes
	var income model.Incomes
	statement, err := repository.DatabaseConnection.Prepare(`SELECT * FROM timesheet.incomes WHERE member_id = ? AND year = ? AND month = ?`)
	if err != nil {
		return nil, err
	}
	row, err := statement.Query(memberID, year, month)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		err = row.Scan(
			&income.ID,
			&income.MemberID,
			&income.Month,
			&income.Year,
			&income.Day,
			&income.StartTimeAMHours,
			&income.StartTimeAMMinutes,
			&income.StartTimeAMSeconds,
			&income.EndTimeAMHours,
			&income.EndTimeAMMinutes,
			&income.EndTimeAMSeconds,
			&income.StartTimePMHours,
			&income.StartTimePMMinutes,
			&income.StartTimePMSeconds,
			&income.EndTimePMHours,
			&income.EndTimePMMinutes,
			&income.EndTimePMSeconds,
			&income.Overtime,
			&income.TotalHoursHours,
			&income.TotalHoursMinutes,
			&income.TotalHoursSeconds,
			&income.CoachingCustomerCharging,
			&income.CoachingPaymentRate,
			&income.TrainingWage,
			&income.OtherWage,
			&income.Company,
			&income.Description,
		)
		if row.Err() != nil {
			return nil, err
		}
		incomeList = append(incomeList, income)
	}
	return incomeList, nil
}

func (repository TimesheetRepository) CreateTransactionTimsheet(transactionTimesheet model.TransactionTimesheet, transactionID string) error {
	statement, err := repository.DatabaseConnection.Prepare(`INSERT INTO transactions (id, member_id, month, year, company, member_name_th, coaching, training, other, total_incomes, salary, income_tax_1, social_security, net_salary, wage, income_tax_53_percentage, income_tax_53, net_wage, net_transfer) VALUES ( ? , ? ,? , ? ,? , ? ,? , ? ,? , ? ,? , ? ,? , ? ,? , ? ,? , ? , ? )`)
	if err != nil {
		return err
	}
	_, err = statement.Exec(
		transactionID,
		transactionTimesheet.MemberID,
		transactionTimesheet.Month,
		transactionTimesheet.Year,
		transactionTimesheet.Company,
		transactionTimesheet.MemberNameTH,
		transactionTimesheet.Coaching,
		transactionTimesheet.Training,
		transactionTimesheet.Other,
		transactionTimesheet.TotalIncomes,
		transactionTimesheet.Salary,
		transactionTimesheet.IncomeTax1,
		transactionTimesheet.SocialSecurity,
		transactionTimesheet.NetSalary,
		transactionTimesheet.Wage,
		transactionTimesheet.IncomeTax53Percentage,
		transactionTimesheet.IncomeTax53,
		transactionTimesheet.NetWage,
		transactionTimesheet.NetTransfer)
	return err
}

func (repository TimesheetRepository) UpdateTransactionTimsheet(transactionTimesheet model.TransactionTimesheet, transactionID string) error {
	statement, err := repository.DatabaseConnection.Prepare(`UPDATE transactions SET coaching = ?, training = ?, other = ?, total_incomes = ?, salary = ?, income_tax_1 = ?, social_security = ?, net_salary = ?, wage = ?, income_tax_53_percentage = ?, income_tax_53 = ?, net_wage = ?, net_transfer = ? WHERE id = ?`)
	if err != nil {
		return err
	}
	_, err = statement.Exec(
		transactionTimesheet.Coaching,
		transactionTimesheet.Training,
		transactionTimesheet.Other,
		transactionTimesheet.TotalIncomes,
		transactionTimesheet.Salary,
		transactionTimesheet.IncomeTax1,
		transactionTimesheet.SocialSecurity,
		transactionTimesheet.NetSalary,
		transactionTimesheet.Wage,
		transactionTimesheet.IncomeTax53Percentage,
		transactionTimesheet.IncomeTax53,
		transactionTimesheet.NetWage,
		transactionTimesheet.NetTransfer,
		transactionID)
	if err != nil {
		return err
	}
	return nil
}

func (repository TimesheetRepository) CreateTimesheet(payment model.Payment, timesheetID, memberID string, year int, month int) error {
	statement, err := repository.DatabaseConnection.Prepare(`INSERT INTO timesheets (id, member_id, month, year, total_hours_hours, total_hours_minutes, total_hours_seconds, total_coaching_customer_charging, total_coaching_payment_rate, total_training_wage, total_other_wage, payment_wage) VALUES ( ? , ? , ? ,? , ? ,? , ? ,? , ? ,? , ? ,? )`)
	if err != nil {
		return err
	}
	_, err = statement.Exec(
		timesheetID,
		memberID,
		month,
		year,
		payment.TotalHoursHours,
		payment.TotalHoursMinutes,
		payment.TotalHoursSeconds,
		payment.TotalCoachingCustomerCharging,
		payment.TotalCoachingPaymentRate,
		payment.TotalTrainigWage,
		payment.TotalOtherWage,
		payment.PaymentWage)
	if err != nil {
		return err
	}
	return nil
}

func (repository TimesheetRepository) UpdateTimesheet(payment model.Payment, timesheetID string) error {
	statement, err := repository.DatabaseConnection.Prepare(`UPDATE timesheets SET total_hours_hours = ?, total_hours_minutes = ?, total_hours_seconds = ?, total_coaching_customer_charging = ?, total_coaching_payment_rate = ?, total_training_wage = ?, total_other_wage = ?, payment_wage = ? WHERE id = ?`)
	if err != nil {
		return err
	}
	_, err = statement.Exec(
		payment.TotalHoursHours,
		payment.TotalHoursMinutes,
		payment.TotalHoursSeconds,
		payment.TotalCoachingCustomerCharging,
		payment.TotalCoachingPaymentRate,
		payment.TotalTrainigWage,
		payment.TotalOtherWage,
		payment.PaymentWage,
		timesheetID)
	if err != nil {
		return err
	}
	return nil
}
