package app

import (
	"math/rand"
	"sort"
	"time"
)

const (
	populationSize         = 9
	turnamentSelectionSize = 3
	numbOfEliteSchedules   = 1
	mutationRate           = 0.2
)

// RunAlgorithm runs the algorithm
func RunAlgorithm(courses []*Course, students []*Student, maxStudents, maxGenerations int) *Schedule {
	Courses = courses
	Students = students
	MaxStudents = maxStudents

	// Set the seeder for our random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	lastFitness := float32(0)
	genNumber := 0
	gensNotChanged := 0

	population := newPopulation(populationSize)

	for gensNotChanged < maxGenerations {
		genNumber++

		sort.Sort(sort.Reverse(ByFitness(population.Schedules)))
		population = evolve(population)

		if population.Schedules[0].Fitness == lastFitness {
			gensNotChanged++
		} else {
			gensNotChanged = 0
		}

		lastFitness = population.Schedules[0].Fitness
	}

	return population.Schedules[0]
}

// Evolve evloves the population to a next gen
func evolve(pop *Population) *Population {
	return mutatePopulation(crossoverPopulation(pop))
}

func crossoverPopulation(pop *Population) *Population {
	crossoverPop := newPopulation(0)

	for i := 0; i < numbOfEliteSchedules; i++ {
		crossoverPop.Schedules = append(crossoverPop.Schedules, pop.Schedules[i])
	}

	for i := numbOfEliteSchedules; i < populationSize; i++ {
		scheduleOne := selectTurnamentPopulation(pop).Schedules[0]
		scheduleTwo := selectTurnamentPopulation(pop).Schedules[0]
		crossoverSchedule := crossoverSchedule(scheduleOne, scheduleTwo)
		crossoverPop.Schedules = append(crossoverPop.Schedules, crossoverSchedule)
	}
	return crossoverPop
}

func mutatePopulation(pop *Population) *Population {
	for i := numbOfEliteSchedules; i < populationSize; i++ {
		mutateSchedule(pop.Schedules[i])
	}
	return pop
}

func mutateSchedule(mutSchedule *Schedule) *Schedule {
	schedule := newSchedule()

	for i := 0; i < len(mutSchedule.Weeks); i++ {
		if mutationRate > rand.Float32() {
			mutSchedule.Weeks[i] = schedule.Weeks[i]
		}
	}

	return mutSchedule
}

func crossoverSchedule(first *Schedule, second *Schedule) *Schedule {
	crossoverSchedule := newSchedule()

	for i := 0; i < len(crossoverSchedule.Weeks); i++ {
		if rand.Intn(2) == 1 {
			crossoverSchedule.Weeks[i] = first.Weeks[i]
		} else {
			crossoverSchedule.Weeks[i] = second.Weeks[i]
		}
	}
	return crossoverSchedule
}

func selectTurnamentPopulation(pop *Population) *Population {
	turnamentPop := newPopulation(0)

	for i := 0; i < turnamentSelectionSize; i++ {
		schedule := pop.Schedules[rand.Intn(populationSize)]
		turnamentPop.Schedules = append(turnamentPop.Schedules, schedule)
	}

	sort.Sort(sort.Reverse(ByFitness(turnamentPop.Schedules)))

	return turnamentPop
}
