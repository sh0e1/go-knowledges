package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sh0e1/go-knowledges/orm/ent/ent"
	"github.com/sh0e1/go-knowledges/orm/ent/ent/car"
	"github.com/sh0e1/go-knowledges/orm/ent/ent/group"
	"github.com/sh0e1/go-knowledges/orm/ent/ent/user"
)

func main() {
	client, err := ent.Open("mysql", "root@tcp(127.0.0.1:3306)/ent?parseTime=true")
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}
	defer client.Close()

	f, err := os.Create("migrate.sql")
	if err != nil {
		log.Fatalf("create migrate file: %v", err)
	}
	defer f.Close()

	ctx := context.Background()
	if err := client.Schema.WriteTo(ctx, f); err != nil {
		log.Fatalf("failed printing schema changes: %v", err)
	}

	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	//if _, err := CreateUser(ctx, client); err != nil {
	//	log.Fatal(err)
	//}

	//a8m, err := QueryUser(ctx, client)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//if _, err := CreateCars(ctx, client); err != nil {
	//	log.Fatal(err)
	//}

	//if err := QueryCars(ctx, a8m); err != nil {
	//	log.Fatal(err)
	//}

	//if err := QueryCarUsers(ctx, a8m); err != nil {
	//	log.Fatal(err)
	//}

	//if err := CreateGraph(ctx, client); err != nil {
	//	log.Fatal(err)
	//}

	if err := QueryGithub(ctx, client); err != nil {
		log.Fatal(err)
	}

	if err := QueryArielCars(ctx, client); err != nil {
		log.Fatal(err)
	}

	if err := QueryGroupWithUsers(ctx, client); err != nil {
		log.Fatal(err)
	}
}

func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.
		Create().
		SetAge(30).
		SetName("a8m").
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}

	log.Println("user was created: ", u)
	return u, nil
}

func QueryUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.
		Query().
		Where(user.Name("a8m")).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %w", err)
	}

	log.Println("user returned: ", u)
	return u, nil
}

func CreateCars(ctx context.Context, client *ent.Client) (*ent.User, error) {
	tesla, err := client.Car.
		Create().
		SetModel("Tesla").
		SetRegisteredAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating car: %w", err)
	}
	log.Println("car was created: ", tesla)

	ford, err := client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating car: %w", err)
	}
	log.Println("car was created: ", ford)

	a8m, err := client.User.
		Create().
		SetAge(30).
		SetName("a8m").
		AddCars(tesla, ford).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", a8m)

	return a8m, nil
}

func QueryCars(ctx context.Context, a8m *ent.User) error {
	cars, err := a8m.QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %w", err)
	}
	log.Println("returned cars: ", cars)

	ford, err := a8m.QueryCars().
		Where(car.Model("Ford")).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %w", err)
	}
	log.Println(ford)

	return nil
}

func QueryCarUsers(ctx context.Context, a8m *ent.User) error {
	cars, err := a8m.QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %w", err)
	}

	for _, ca := range cars {
		owner, err := ca.QueryOwner().Only(ctx)
		if err != nil {
			return fmt.Errorf("failed querying car %q woner: %w", ca.Model, err)
		}
		log.Printf("car %q owner: %q\n", ca.Model, owner.Name)
	}

	return nil
}

func CreateGraph(ctx context.Context, client *ent.Client) error {
	a8m, err := client.User.
		Create().
		SetAge(30).
		SetName("Ariel").
		Save(ctx)
	if err != nil {
		return err
	}

	neta, err := client.User.
		Create().
		SetAge(28).
		SetName("Neta").
		Save(ctx)
	if err != nil {
		return err
	}

	if _, err := client.Car.
		Create().
		SetModel("Tesla").
		SetRegisteredAt(time.Now()).
		SetOwner(a8m).
		Save(ctx); err != nil {
		return err
	}

	if _, err := client.Car.
		Create().
		SetModel("Mazda").
		SetRegisteredAt(time.Now()).
		SetOwner(a8m).
		Save(ctx); err != nil {
		return err
	}
	if _, err := client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Now()).
		SetOwner(neta).
		Save(ctx); err != nil {
		return err
	}

	if _, err = client.Group.
		Create().
		SetName("GitLab").
		AddUsers(neta, a8m).
		Save(ctx); err != nil {
		return err
	}

	if _, err := client.Group.
		Create().
		SetName("GitHub").
		AddUsers(a8m).
		Save(ctx); err != nil {
		return err
	}

	log.Println("The graph was created successfully")
	return nil
}

func QueryGithub(ctx context.Context, client *ent.Client) error {
	cars, err := client.Group.
		Query().
		Where(group.Name("GitHub")).
		QueryUsers().
		QueryCars().
		All(ctx)
	if err != nil {
		return fmt.Errorf("failed getting cars: %w", err)
	}

	log.Println("cars returned:", cars)
	return nil
}

func QueryArielCars(ctx context.Context, client *ent.Client) error {
	a8m := client.User.
		Query().
		Where(
			user.HasCars(),
			user.Name("Ariel"),
		).
		OnlyX(ctx)
	cars, err := a8m.
		QueryGroups().
		QueryUsers().
		QueryCars().
		Where(
			car.Not(
				car.Model("Mazda"),
			),
		).
		All(ctx)
	if err != nil {
		return fmt.Errorf("failed getting cars: %w", err)
	}

	log.Println("cars returned:", cars)
	return nil
}

func QueryGroupWithUsers(ctx context.Context, client *ent.Client) error {
	groups, err := client.Group.
		Query().
		Where(group.HasUsers()).
		All(ctx)
	if err != nil {
		return fmt.Errorf("failed getting groups: %w", err)
	}

	log.Println("groups returned:", groups)
	return nil
}
