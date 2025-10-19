package utils

// tree_util.go

import "fmt"

// TreeNode 通用树节点结构
type TreeNode[K comparable, E any] struct {
	ID       K                 `json:"id"`        // 主键
	ParentID K                 `json:"parent_id"` // 父节点ID
	Name     string            `json:"name"`      // 名称
	Sort     int               `json:"sort"`      // 排序权重
	Deep     int               `json:"deep"`      // 层级深度（从 0 开始）
	Extra    E                 `json:"extra"`     // 额外数据（原始对象）
	Children []*TreeNode[K, E] `json:"children"`  // 子节点列表
}

// ListToTree 将扁平列表转换为树结构
// 参数：
//
//	list: 扁平的节点切片
//	rootParentID: 根节点的 ParentID（例如 0 或 ""）
//
// 返回：
//
//	根节点列表（森林），error
func ListToTree[K comparable, T any](
	list []TreeNode[K, T],
	rootParentID K,
) ([]*TreeNode[K, T], error) {
	nodeMap := make(map[K]*TreeNode[K, T])
	var roots []*TreeNode[K, T]

	// 第一步：将所有节点放入 map，方便查找
	for i := range list {
		node := &list[i]
		nodeMap[node.ID] = node
		node.Children = []*TreeNode[K, T]{} // 初始化
	}

	// 第二步：建立父子关系
	for i := range list {
		node := &list[i]

		if node.ParentID == rootParentID || node.ParentID == node.ID {
			// 是根节点，或指向自己
			roots = append(roots, node)
		} else {
			parentNode, exists := nodeMap[node.ParentID]
			if !exists {
				return nil, fmt.Errorf("parent node not found for node ID=%v, ParentID=%v", node.ID, node.ParentID)
			}
			parentNode.Children = append(parentNode.Children, node)
		}

		// 设置深度（可选）
		if node.ParentID == rootParentID {
			node.Deep = 0
		} else {
			parentNode, exists := nodeMap[node.ParentID]
			if exists {
				node.Deep = parentNode.Deep + 1
			} else {
				node.Deep = 0 // 默认
			}
		}
	}

	return roots, nil
}

// TreeToList 将树结构展开为扁平列表（前序遍历）
func TreeToList[K comparable, T any](roots []*TreeNode[K, T]) []TreeNode[K, T] {
	var result []TreeNode[K, T]

	var dfs func(node *TreeNode[K, T])
	dfs = func(node *TreeNode[K, T]) {
		if node == nil {
			return
		}
		// 复制节点（避免指针问题）
		copied := *node
		copied.Children = nil // 不包含子节点
		result = append(result, copied)

		for _, child := range node.Children {
			dfs(child)
		}
	}

	for _, root := range roots {
		dfs(root)
	}

	return result
}
