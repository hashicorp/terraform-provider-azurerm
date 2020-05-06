package tests

import AzureRM
import org.junit.Assert.assertTrue
import org.junit.Test

class VcsTests {
    @Test
    fun buildsHaveCleanCheckOut() {
        val project = AzureRM("public", TestConfiguration())
        project.buildTypes.forEach { bt ->
            assertTrue("Build '${bt.id}' doesn't use clean checkout", bt.vcs.cleanCheckout)
        }
    }
}
