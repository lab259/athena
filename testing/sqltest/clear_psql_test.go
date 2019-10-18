package sqltest_test

import (
	"database/sql"

	"github.com/lab259/athena/testing/sqltest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	_ "github.com/lib/pq"
)

var _ = Describe("ClearPostgreSQL", func() {
	var db *sql.DB

	BeforeEach(func() {
		var err error
		db, err = sql.Open("postgres", `user=postgres dbname=postgres sslmode=disable`)
		Expect(err).ToNot(HaveOccurred())
		Expect(db.Ping()).To(Succeed())

		_, err = db.Exec(`
		drop table if exists books;

		drop table if exists movies;

		drop table if exists songs;

		create table books (
			id serial primary key,
			name varchar (50)
		);

		create table movies (
			id serial primary key,
			name varchar (50)
		);

		create table songs (
			id serial primary key,
			name varchar (50)
		);

		insert into books (name)
		values
		    ('The Great Gatsby'),
			('To Kill a Mockingbird'),
			('1984');

		insert into movies (name)
		values 
			('Forrest Gump'),
			('Bonnie and Clyde'),
			('2001: A Space Odyssey'),
			('Star Wars');
		
		insert into songs (name)
		values 
			('Thriller'),
			('Like a Prayer');
		`)

		Expect(err).ToNot(HaveOccurred())

		var count int64
		Expect(db.QueryRow(`select count(*) from books`).Scan(&count)).To(Succeed())
		Expect(count).To(BeNumerically("==", 3))

		Expect(db.QueryRow(`select count(*) from movies`).Scan(&count)).To(Succeed())
		Expect(count).To(BeNumerically("==", 4))

		Expect(db.QueryRow(`select count(*) from songs`).Scan(&count)).To(Succeed())
		Expect(count).To(BeNumerically("==", 2))
	})

	It("should clear all tables", func() {
		sqltest.ClearPostgreSQL(1, db)

		var count int64
		Expect(db.QueryRow(`select count(*) from books`).Scan(&count)).To(Succeed())
		Expect(count).To(BeNumerically("==", 0))

		Expect(db.QueryRow(`select count(*) from movies`).Scan(&count)).To(Succeed())
		Expect(count).To(BeNumerically("==", 0))

		Expect(db.QueryRow(`select count(*) from songs`).Scan(&count)).To(Succeed())
		Expect(count).To(BeNumerically("==", 0))
	})

	It("should ignore specified table", func() {
		sqltest.ClearPostgreSQL(1, db, "movies")

		var count int64
		Expect(db.QueryRow(`select count(*) from books`).Scan(&count)).To(Succeed())
		Expect(count).To(BeNumerically("==", 0))

		Expect(db.QueryRow(`select count(*) from movies`).Scan(&count)).To(Succeed())
		Expect(count).To(BeNumerically("==", 4))

		Expect(db.QueryRow(`select count(*) from songs`).Scan(&count)).To(Succeed())
		Expect(count).To(BeNumerically("==", 0))
	})

	It("should ignore specified tables", func() {
		sqltest.ClearPostgreSQL(1, db, "movies", "books")

		var count int64
		Expect(db.QueryRow(`select count(*) from books`).Scan(&count)).To(Succeed())
		Expect(count).To(BeNumerically("==", 3))

		Expect(db.QueryRow(`select count(*) from movies`).Scan(&count)).To(Succeed())
		Expect(count).To(BeNumerically("==", 4))

		Expect(db.QueryRow(`select count(*) from songs`).Scan(&count)).To(Succeed())
		Expect(count).To(BeNumerically("==", 0))
	})

})
