package integrationTest

import (
	"database/sql"
	"fmt"
	"os/exec"
	"reflect"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("user alias", func() {
	BeforeEach(func() {
		err := resetDB()
		Expect(err).NotTo(HaveOccurred())
	})

	It("can add an user alias", func() {
		if skipUserAlias && !isCI {
			Skip("can add an user alias")
		}

		cli := exec.Command(cliPath, "user", "add", userName1, userPW)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "user", "add-alias", alias1, userName1)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := fmt.Sprintf("Successfully added user alias %v -> %v\n", alias1, userName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool
		domain := strings.Split(userName1, "@")[1]

		query := `SELECT exists
		(SELECT * FROM forwardings
		WHERE address = '` + fmt.Sprintf("%v@%v", alias1, domain) + `' AND forwarding = '` + userName1 + `'
		AND is_alias = 1 AND active = 1 AND is_forwarding = 0 AND is_maillist = 0);`

		err = db.QueryRow(query).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))
	})

	It("can't add an user alias if email exists", func() {
		if skipUserAlias && !isCI {
			Skip("can't add an user alias if email exists")
		}

		cli := exec.Command(cliPath, "user", "add", userName1, userPW)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		name := strings.Split(userName1, "@")[0]

		cli = exec.Command(cliPath, "user", "add-alias", name, userName1)
		output, err := cli.CombinedOutput()
		Expect(err).To(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("An user with %v already exists\n", userName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

	It("can't add an user alias if user alias already exists", func() {
		if skipUserAlias && !isCI {
			Skip("can't add an user alias if user alias already exists")
		}

		cli := exec.Command(cliPath, "user", "add", userName1, userPW)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "user", "add-alias", alias1, userName1)
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "user", "add-alias", alias1, userName1)
		output, err := cli.CombinedOutput()
		Expect(err).To(HaveOccurred())

		domain := strings.Split(userName1, "@")[1]

		actual := string(output)
		expected := fmt.Sprintf("An alias with %v already exists\n", fmt.Sprintf("%v@%v", alias1, domain))

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

	It("can delete an user alias", func() {
		if skipUserAlias && !isCI {
			Skip("can delete an user alias")
		}

		cli := exec.Command(cliPath, "user", "add", userName1, userPW)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "user", "add-alias", alias1, userName1)
		err = cli.Run()
		Expect(err).NotTo(HaveOccurred())

		domain := strings.Split(userName1, "@")[1]

		cli = exec.Command(cliPath, "user", "delete-alias", fmt.Sprintf("%v@%v", alias1, domain))
		output, err := cli.CombinedOutput()
		Expect(err).ToNot(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("Successfully deleted user alias %v\n", fmt.Sprintf("%v@%v", alias1, domain))

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

	It("can't delete an alias which doesn't exist", func() {
		if skipUserAlias && !isCI {
			Skip("can't delete an alias which doesn't exist")
		}

		cli := exec.Command(cliPath, "user", "delete-alias", userName1)
		output, err := cli.CombinedOutput()
		Expect(err).To(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("An alias with %v doesn't exists\n", userName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

})