package main

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"
)

func placeCityCrawlerAtWork(s *State, city City) {
	stage, building := cityNextBuildingStage(s, city)
	if !building {
		return
	}
	id := cityCrawlerID(s, city)
	crawler := s.Enemies[id]
	crawler.X, crawler.Y = cityCrawlerWorkPosition(city, stage)
	s.Enemies[id] = crawler
}

func TestCityCrawlerVisitsEveryBuildingBeforeWorking(t *testing.T) {
	s := newGame()
	noRivals(s)
	cityID := s.foundCity(500, 500, 0.4)
	crawlerID := cityCrawlerID(s, s.Cities[cityID])
	for stage := range cityBuildOrder {
		city := s.Cities[cityID]
		crawler := s.Enemies[crawlerID]
		x, y := cityCrawlerWorkPosition(city, stage)
		if math.Hypot(crawler.X-x, crawler.Y-y) < cityCrawlerSpeed/60 {
			t.Fatal("successive city jobs do not require travel")
		}
		runTicks(s, 1)
		moved := s.Enemies[crawlerID]
		gap := math.Hypot(moved.X-crawler.X, moved.Y-crawler.Y)
		if math.Abs(gap-cityCrawlerSpeed/60) > 1e-9 {
			t.Fatalf("constructor moved %.3f m, want one driving step", gap)
		}
		if moved.Facing != facingFromMovement(
			moved.X-crawler.X, moved.Y-crawler.Y,
		) || s.Cities[cityID].Work != cityBuildTicks {
			t.Fatal("travel did not turn the crawler or advanced building work")
		}
		if !tickUntil(s, 30*60, func() bool {
			return s.Cities[cityID].Work < cityBuildTicks
		}) {
			t.Fatal("the constructor never reached its next job")
		}
		crawler = s.Enemies[crawlerID]
		if math.Hypot(crawler.X-x, crawler.Y-y) > 1e-9 {
			t.Fatal("construction started before the crawler arrived")
		}
		bx, by := cityBuildingPosition(city, stage)
		if math.Hypot(crawler.X-bx, crawler.Y-by) >= cityCrawlerBubbleUnits {
			t.Fatal("the constructor's bubble does not reach its building")
		}
		runTicks(s, cityBuildTicks-2)
		if s.Cities[cityID].Stage != stage ||
			s.Cities[cityID].Work != 1 {
			t.Fatal("the building finished before 45 seconds of on-site work")
		}
		runTicks(s, 1)
		if s.Cities[cityID].Stage != stage+1 ||
			s.Enemies[crawlerID].Health != enemySpecOf(EnemyCityCrawler).health {
			t.Fatal("the building did not finish or its constructor lost hull")
		}
	}
}

func TestCityCrawlerTravelSurvivesSaving(t *testing.T) {
	s := newGame()
	noRivals(s)
	id := s.foundCity(500, 500, 0.7)
	runTicks(s, 100)
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var loaded State
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatal(err)
	}
	loaded.enterRegion()
	runTicks(s, 1000)
	runTicks(&loaded, 1000)
	if !reflect.DeepEqual(s.Cities[id], loaded.Cities[id]) ||
		!reflect.DeepEqual(s.Enemies, loaded.Enemies) {
		t.Fatal("saving changed the constructor's route or construction work")
	}
}

func TestCityCrawlerDrivesBackToDestroyedBuildings(t *testing.T) {
	s := newGame()
	noRivals(s)
	id := finishedCityForTest(s)
	city := s.Cities[id]
	crawlerID := cityCrawlerID(s, city)
	crawler := s.Enemies[crawlerID]
	crawler.X, crawler.Y = cityCrawlerWorkPosition(city, 4)
	s.Enemies[crawlerID] = crawler
	var oilworksID int64
	for _, buildingID := range city.BuildingIDs {
		if s.Enemies[buildingID].Kind == EnemyCityOilworks {
			oilworksID = buildingID
		}
	}
	s.killEnemy(oilworksID)
	runTicks(s, 1)
	if s.Cities[id].Work != cityBuildTicks {
		t.Fatal("reconstruction advanced before the constructor arrived")
	}
	if !tickUntil(s, 30*60, func() bool {
		return s.Cities[id].Work < cityBuildTicks
	}) {
		t.Fatal("the constructor did not drive back to the destroyed oilworks")
	}
	x, y := cityCrawlerWorkPosition(city, 2)
	crawler = s.Enemies[crawlerID]
	if math.Hypot(crawler.X-x, crawler.Y-y) > 1e-9 {
		t.Fatal("the constructor started rebuilding away from the oilworks")
	}
}
