package txnnode

import "context"

func createWriteIntent(
	ctx context.Context,
) {

	//TODO get newer committed value
	// if exists, restart

	//TODO check write intent
	// if exists, report conflict

	//TODO check timestamp cache
	// if older, restart

}
