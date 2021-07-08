#!/bin/bash
docker-compose -f docker-compose.test.yml up -d

is_finished() {
	local service_name="$1"
	local status="$(docker inspect -f "{{.State.Status}}" "$service_name")"

	echo "STATUS: $status, CONTAINER: $service_name"

	if [ "$status" = "exited" ]; then
		return 0
	else
		return 1
	fi
}

return_exit_code() {
	local service_name="$1"
	local exit_code="$(docker inspect -f "{{.State.ExitCode}}" "$service_name")"
	
	echo "EXIT_CODE: $exit_code, CONTAINER: $service_name"

	if [ "$exit_code" = 0 ]; then
		return 0
	else
		return 1
	fi
}

while ! is_finished nistagram-posts-test; do sleep 20; done

NISTAGRAM_POSTS_TEST_EXIT_CODE="$(return_exit_code nistagram-posts-test)"

echo "nistagram-posts tests returned $NISTAGRAM_POSTS_TEST_EXIT_CODE"

if [ "$NISTAGRAM_POSTS_TEST_EXIT_CODE" -eq 1 ]; then
	echo "::set-output name=tests_exit_code::$($NISTAGRAM_POSTS_TEST_EXIT_CODE)"
fi

echo "::set-output name=tests_exit_code::0"