package impl

const (
	insertBook = `INSERT INTO books(id,status,create_at,create_by,name,author) VALUES (?,?,?,?,?,?);`

	queryBook = `SELECT * FROM books Where name Like ? AND author Like ? LIMIT ?,?`

	updateBook = `UPDATE books SET update_at=?,update_by=?,name=?,author=? WHERE id =?`

	deleteBook = `DELETE FROM books WHERE id = ?`
)
