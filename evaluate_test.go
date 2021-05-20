package typescript

import (
	"context"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var (
	amdModuleScript = strings.TrimSpace(`
		define("myModule", ["require", "exports"], function (require, exports, core_1) {
			Object.defineProperty(exports, "__esModule", { value: true });
			exports.multiply = void 0;
			function multiply() {
				var nums = [];
				for (var _i = 0; _i < arguments.length; _i++) {
					nums[_i] = arguments[_i];
				}
				if (Array.isArray(nums)) {
					if (nums.length === 0) {
						return 0;
					}
					else if (nums.length === 1) {
						return nums[0];
					}
				}
				else {
					return 0;
				}
				var a = nums[0];
				for (var i = 1; i < nums.length; i++) {
					a += a * nums[i];
				}
				return a;
			}
			exports.multiply = multiply;
		});
	`)
)

func TestEvaluateCtx(t *testing.T) {
	// This test hits a lot of things:
	//  #1 - We test that we can load the almond AMD module loader
	//  #2 - We test that we can load our own 'evaluate before' script that declares an AMD module
	//  #3 - We test that we can import the AMD module in a type script script and use a function from the module
	t.Run("evaluate with custom AMD module", func(t *testing.T) {
		script := "import { multiply } from 'myModule'; multiply(5, 5)"
		result, err := EvaluateCtx(context.Background(), strings.NewReader(script),
			WithAlmondModuleLoader(),
			WithTranspile(true),
			WithEvaluateBefore(strings.NewReader(amdModuleScript)),
			WithTranspileOptions(WithVerbose()))
		require.NoError(t, err)
		require.Equal(t, int64(30), result.ToInteger())
	})
}
