package restapi

import "os"

func SetupEnv() {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "josh")
	os.Setenv("DB_PASSWORD", "111kkk")
	os.Setenv("DB_NAME", "upsize")
	os.Setenv("TEST_DB_HOST", "localhost")
	os.Setenv("TEST_DB_PORT", "5432")
	os.Setenv("TEST_DB_USER", "josh")
	os.Setenv("TEST_DB_PASSWORD", "111kkk")
	os.Setenv("TEST_DB_NAME", "upsize_test")
}
