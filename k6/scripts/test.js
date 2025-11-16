import http from 'k6/http'
import { check } from 'k6'

const BASE_URL = __ENV.BASE_URL || 'http://api:8080'

export const options = {
  scenarios: {
    reviewer_service: {
      executor: 'constant-arrival-rate',
      rate: 5,
      timeUnit: '1s',
      duration: '1m',
      preAllocatedVUs: 10,
      maxVUs: 50,
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<300'],
    checks: ['rate>0.999'],
  },
}

function makeName(prefix) {
  return `${prefix}-${__VU}-${Date.now()}-${__ITER}`
}

function jsonHeaders() {
  return { headers: { 'Content-Type': 'application/json' } }
}

export default function () {
  const teamName = makeName('team')
  const authorID = `${teamName}-author`
  const reviewer1 = `${teamName}-rev1`
  const reviewer2 = `${teamName}-rev2`
  const reviewer3 = `${teamName}-rev3`

  const teamPayload = {
    team_name: teamName,
    members: [
      { user_id: authorID, username: `${teamName}-author`, is_active: true },
      { user_id: reviewer1, username: `${teamName}-rev1`, is_active: true },
      { user_id: reviewer2, username: `${teamName}-rev2`, is_active: true },
      { user_id: reviewer3, username: `${teamName}-rev3`, is_active: true },
    ],
  }

  const createTeamResp = http.post(`${BASE_URL}/team/add`, JSON.stringify(teamPayload), jsonHeaders())
  check(createTeamResp, {
    'team/add success': (r) => r.status === 201,
  })

  const getTeamResp = http.get(`${BASE_URL}/team/get?team_name=${teamName}`)
  check(getTeamResp, {
    'team/get success': (r) => r.status === 200,
  })

  const prID = makeName('pr')
  const createPrResp = http.post(
    `${BASE_URL}/pullRequest/create`,
    JSON.stringify({
      pull_request_id: prID,
      pull_request_name: `Load test ${prID}`,
      author_id: authorID,
    }),
    jsonHeaders(),
  )

  check(createPrResp, {
    'pullRequest/create success': (r) => r.status === 201,
  })

  const prBody = createPrResp.json('pr')
  const assigned = prBody?.assigned_reviewers || []
  const oldReviewer = assigned[0] || reviewer1

  const getReviewResp = http.get(`${BASE_URL}/users/getReview?user_id=${oldReviewer}`)
  check(getReviewResp, {
    'users/getReview success': (r) => r.status === 200,
  })

  const reassignResp = http.post(
    `${BASE_URL}/pullRequest/reassign`,
    JSON.stringify({
      pull_request_id: prID,
      old_user_id: oldReviewer,
    }),
    jsonHeaders(),
  )

  check(reassignResp, {
    'pullRequest/reassign success': (r) => r.status === 200,
  })

  const mergeResp = http.post(
    `${BASE_URL}/pullRequest/merge`,
    JSON.stringify({
      pull_request_id: prID,
    }),
    jsonHeaders(),
  )

  check(mergeResp, {
    'pullRequest/merge success': (r) => r.status === 200,
  })
}
