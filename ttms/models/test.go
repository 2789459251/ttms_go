package models

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

var cache []int
var i int = 0
var tmp int

func isSymmetric(root *TreeNode) bool {
	Cache(root)
	return check(root)
}
func Cache(root *TreeNode) {
	if root == nil {
		return
	}
	cache = append(cache, root.Val)
	Cache(root.Left)
	Cache(root.Right)
}
func check(root *TreeNode) bool {
	if i == len(cache) && root == nil {
		return true
	}
	tmp = i
	i++
	if root == nil || i == len(cache) || root.Val != cache[tmp] {
		return false
	}
	return check(root.Right) && check(root.Left)
}
