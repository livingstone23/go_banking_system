package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)


//New interface
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)


} 




// Store provides all functions to execute db queries and transactions
type SQLStore struct {
	*Queries // this is the queries struct
	//db *pgx.Conn 		// this is the connection to the database
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {

	txOptions := pgx.TxOptions{} // Define las opciones de transacción aquí

	//tx, err := store.db.BeginTx(ctx, nil)	//TODO: check if this is the correct way to start a transaction
	tx, err := store.db.BeginTx(ctx, txOptions) // Pasa txOptions en lugar de nil
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// txKey is an identifier for the transaction
//var txKey = struct{}{}

// TranferTx performs a money transfer from one account to the other
// It creates a transfer record, add account entries, and update account's balance within a single database transaction
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		//txName := ctx.Value(txKey)

		//fmt.Println(txName, "Create the transfer...")

		// This is the first step of the transaction
		// It creates a transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		//fmt.Println(txName, "Create entry 1...")
		// This is the second step of the transaction
		// It creates an account entry for the account that the money is transferred from
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		//fmt.Println(txName, "Create entry 2...")
		// This is the third step of the transaction
		// It creates an account entry for the account that the money is transferred to
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		/*
		//fmt.Println(txName, "Get Account 1...")
		// This is the fourth step of the transaction
		// It updates the account balance for the account that the money is transferred from
		//TODO update the account balance for the account that the money is transferred from
		
			Account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
			if err != nil {
				return err
			}

			//fmt.Println(txName, "Update Account 1...")
			result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      arg.FromAccountID,
				Balance: Account1.Balance - arg.Amount,
			})
			if err != nil {
				return err
			}

			//fmt.Println(txName, "Get Account 2...")
			Account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
			if err != nil {
				return err
			}

			//fmt.Println(txName, "Update Account 2...")
			result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      arg.ToAccountID,
				Balance: Account2.Balance + arg.Amount,
			})
			if err != nil {
				return err
			}
		*/

		if arg.FromAccountID < arg.ToAccountID {

			/*
			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.FromAccountID,
				Amount: -arg.Amount,
			})
			if err != nil {
				return err
			}

			//fmt.Println(txName, "Update Account 2...")
			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount,
			})
			if err != nil {
				return err
			}
			*/
			
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)

		} else {

			/*
			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount,
			})
			if err != nil {
				return err
			}

			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.FromAccountID,
				Amount: -arg.Amount,
			})
			if err != nil {
				return err
			}
			*/

			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)

		}

		// This is the fifth step of the transaction
		return nil
	})

	return result, err
}


// TransferTx performs a money transfer from one account to the other
func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return
	}
	return
}

/*
// NewStore creates a new store
func NewStore(db *pgx.Conn) *Store {
	return &Store{
		db: db,
		Queries: New(db),

	}
}
*/
